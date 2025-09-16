package client

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/magiconair/properties"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

// ConfigParser 配置解析器实例
// 核心能力：
// 1. 支持6种配置格式：JSON/YAML/TOML/INI/Properties/Text
// 2. 提供类型安全的配置获取（字符串/数值/布尔/时间/切片/Map等）
// 3. 高并发优化：读写锁保护、缓存机制（默认启用）
// 4. 动态更新：支持运行时更新配置内容
type ConfigParser struct {
	content      string                 // 原始配置内容字符串
	format       string                 // 配置格式（json/yaml/toml/ini/properties/text）
	data         map[string]interface{} // 解析后的数据（展平格式，便于路径访问）
	isTextType   bool                   // 是否为Text格式（仅支持content路径）
	arrayRegex   *regexp.Regexp         // 预编译正则：匹配数组路径（如 "arr[0]"）
	cache        sync.Map               // 缓存已获取的配置值（key：路径，value：配置值）
	cacheEnabled bool                   // 是否启用缓存（默认true）
	mu           sync.RWMutex           // 读写锁：保护共享资源（data/cache/cacheEnabled/isTextType）
}

// NewConfigParser 创建配置解析器实例
// configString：配置内容字符串
// format：配置格式（支持：json/yaml/yml/toml/ini/properties/text，不区分大小写）
// 返回：解析器实例，若格式不支持或解析失败返回nil和错误
func NewConfigParser(configString, format string) (*ConfigParser, error) {
	// 统一格式为小写，避免大小写判断问题
	format = strings.ToLower(format)

	parser := &ConfigParser{
		content:      configString,
		format:       format,
		data:         make(map[string]interface{}),
		arrayRegex:   regexp.MustCompile(`^(\w+)\[(\d+)]$`), // 匹配 "键[索引]" 格式
		cacheEnabled: true,                                  // 默认启用缓存
	}

	switch format {
	case "json":
		if err := parser.parseJSON(); err != nil {
			return nil, fmt.Errorf("json解析失败: %w", err)
		}
	case "yaml", "yml":
		if err := parser.parseYAML(); err != nil {
			return nil, fmt.Errorf("yaml解析失败: %w", err)
		}
	case "toml":
		if err := parser.parseTOML(); err != nil {
			return nil, fmt.Errorf("toml解析失败: %w", err)
		}
	case "ini":
		if err := parser.parseINI(); err != nil {
			return nil, fmt.Errorf("ini解析失败: %w", err)
		}
	case "properties":
		if err := parser.parseProperties(); err != nil {
			return nil, fmt.Errorf("properties解析失败: %w", err)
		}
	case "text":
		parser.data["content"] = configString
		parser.isTextType = true
	default:
		return nil, fmt.Errorf("不支持的配置格式: %s，支持格式：json/yaml/yml/toml/ini/properties/text", format)
	}

	// 非Text格式：展平嵌套结构（如 "a.b.c" 格式），便于路径访问
	if !parser.isTextType {
		parser.data = parser.flattenData(parser.data)
	}

	return parser, nil
}

// GetRawContent 获取原始配置内容字符串
func (p *ConfigParser) GetRawContent() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.content
}

// GetFormat 获取配置格式
func (p *ConfigParser) GetFormat() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.format
}

// EnableCache 启用/禁用缓存
// enabled：true=启用，false=禁用（禁用时会清空现有缓存）
func (p *ConfigParser) EnableCache(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cacheEnabled = enabled
	if !enabled {
		p.cache = sync.Map{} // 高效清空缓存（替换原实例，旧实例由GC回收）
	}
}

// InvalidateCache 清空指定路径的缓存，或清空全部缓存
// paths：可变参数，传入路径则清空指定路径缓存；不传入则清空全部缓存
// 注意：会同时清空「值缓存」和「路径存在性缓存」
func (p *ConfigParser) InvalidateCache(paths ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(paths) == 0 {
		p.cache = sync.Map{} // 清空全部缓存
		return
	}

	// 清空指定路径的缓存（值缓存 + 存在性缓存）
	for _, path := range paths {
		p.cache.Delete(path)
		p.cache.Delete("exists:" + path)
	}
}

// IsSet 检查指定路径是否存在有效配置（路径存在且值不为null）
// path：配置路径（如 "a.b.c"、"arr[0]"）
// 返回：true=路径存在且值有效，false=路径不存在或值为null
func (p *ConfigParser) IsSet(path string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Text格式仅支持 "content" 路径
	if p.isTextType {
		return path == "content"
	}

	// 先查存在性缓存（key："exists:路径"）
	existCacheKey := "exists:" + path
	if p.cacheEnabled {
		if cached, found := p.cache.Load(existCacheKey); found {
			return cached.(bool)
		}
	}

	// 实际检查路径（区分「路径不存在」和「值为null」）
	value, err := p.findNestedValue(p.data, strings.Split(path, "."))
	exists := err == nil && value != nil

	// 缓存存在性结果
	if p.cacheEnabled {
		p.cache.Store(existCacheKey, exists)
	}

	return exists
}

// UpdateConfig 动态更新配置内容
// configString：新的配置内容字符串
// format：新配置的格式（需与初始化格式一致，避免类型不兼容）
// 返回：更新成功返回nil，失败返回错误（如解析失败）
// 注意：更新后会清空现有缓存，同步更新Text格式标识
func (p *ConfigParser) UpdateConfig(configString, format string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 先创建新解析器，确保新配置合法（避免更新过程中破坏原配置）
	newParser, err := NewConfigParser(configString, format)
	if err != nil {
		return fmt.Errorf("更新配置失败: %w", err)
	}

	// 同步核心状态（内容 + 格式 + 数据 + Text格式标识）
	p.content = newParser.content
	p.format = newParser.format
	p.data = newParser.data
	p.isTextType = newParser.isTextType

	// 清空缓存（避免旧值干扰）
	p.cache = sync.Map{}

	return nil
}

// ------------------------------
// 核心配置获取方法（无错误返回，失败返回默认值）
// ------------------------------

// GetString 获取指定路径的字符串值
// path：配置路径（如 "a.b.c"、"content"（Text格式））
// 返回：路径对应的字符串值；若路径不存在、转换失败或Text格式路径错误，返回空串
func (p *ConfigParser) GetString(path string) string {
	// Text格式特殊处理：仅支持 "content" 路径
	if p.isTextType {
		p.mu.RLock()
		defer p.mu.RUnlock()

		if path != "content" {
			return ""
		}
		content, ok := p.data["content"].(string)
		if !ok {
			return ""
		}
		return content
	}
	// 非Text格式：获取值并转换为字符串
	value := p.getValueByPathOrNil(path)
	return p.convertToString(value)
}

// GetInt 获取指定路径的int类型值
// path：配置路径（如 "a.b.port"）
// 返回：路径对应的int值；若路径不存在、转换失败或为Text格式，返回0
func (p *ConfigParser) GetInt(path string) int {
	// Text格式不支持数值类型
	if p.isTextType {
		return 0
	}

	// 依赖GetInt64转换，失败返回0
	val, _ := p.toInt64(p.getValueByPathOrNil(path))
	return int(val)
}

// GetInt64 获取指定路径的int64类型值
// path：配置路径（如 "a.b.maxSize"）
// 返回：路径对应的int64值；若路径不存在、转换失败或为Text格式，返回0
func (p *ConfigParser) GetInt64(path string) int64 {
	// Text格式不支持数值类型
	if p.isTextType {
		return 0
	}

	// 获取值并转换，失败返回0
	val, _ := p.toInt64(p.getValueByPathOrNil(path))
	return val
}

// GetUint 获取指定路径的uint类型值
// path：配置路径（如 "a.b.count"）
// 返回：路径对应的uint值；若路径不存在、转换失败或为Text格式，返回0
func (p *ConfigParser) GetUint(path string) uint {
	// Text格式不支持数值类型
	if p.isTextType {
		return 0
	}

	// 依赖GetUint64转换，失败返回0
	val, _ := p.toUint64(p.getValueByPathOrNil(path))
	return uint(val)
}

// GetUint64 获取指定路径的uint64类型值
// path：配置路径（如 "a.b.total"）
// 返回：路径对应的uint64值；若路径不存在、转换失败或为Text格式，返回0
func (p *ConfigParser) GetUint64(path string) uint64 {
	// Text格式不支持数值类型
	if p.isTextType {
		return 0
	}

	// 获取值并转换，失败返回0
	val, _ := p.toUint64(p.getValueByPathOrNil(path))
	return val
}

// GetFloat64 获取指定路径的float64类型值
// path：配置路径（如 "a.b.version"、"a.b.score"）
// 返回：路径对应的float64值；若路径不存在、转换失败或为Text格式，返回0.0
func (p *ConfigParser) GetFloat64(path string) float64 {
	// Text格式不支持数值类型
	if p.isTextType {
		return 0.0
	}

	// 获取值并转换，失败返回0.0
	val, _ := p.toFloat64(p.getValueByPathOrNil(path))
	return val
}

// GetBool 获取指定路径的bool类型值
// path：配置路径（如 "a.b.enable"）
// 返回：路径对应的bool值；若路径不存在、转换失败或为Text格式，返回false
func (p *ConfigParser) GetBool(path string) bool {
	// Text格式不支持布尔类型
	if p.isTextType {
		return false
	}

	// 获取值并转换，失败返回false
	val, _ := p.toBool(p.getValueByPathOrNil(path))
	return val
}

// GetTime 获取指定路径的时间值（time.Time）
// path：配置路径（如 "a.b.startTime"）
// 支持的时间格式：RFC3339、RFC3339Nano、2006-01-02 15:04:05、2006-01-02、15:04:05等
// 返回：路径对应的时间值；若路径不存在、转换失败或为Text格式，返回时间零值（time.Time{}）
func (p *ConfigParser) GetTime(path string) time.Time {
	// Text格式不支持时间类型
	if p.isTextType {
		return time.Time{}
	}

	// 获取值并转换，失败返回时间零值
	val, _ := p.toTime(p.getValueByPathOrNil(path))
	return val
}

// GetDuration 获取指定路径的时间间隔值（time.Duration）
// path：配置路径（如 "a.b.timeout"）
// 支持的格式："30s"、"5m"、"1h" 等（符合time.ParseDuration规范）
// 返回：路径对应的时间间隔值；若路径不存在、转换失败或为Text格式，返回0
func (p *ConfigParser) GetDuration(path string) time.Duration {
	// Text格式不支持时间间隔类型
	if p.isTextType {
		return 0
	}

	// 获取值并转换，失败返回0
	val, _ := p.toDuration(p.getValueByPathOrNil(path))
	return val
}

// GetMap 获取指定路径的map[string]interface{}类型值
// path：配置路径（如 "a.b.params"）
// 返回：路径对应的Map值；若路径不存在、转换失败或为Text格式，返回空Map（非nil）
func (p *ConfigParser) GetMap(path string) map[string]interface{} {
	// Text格式不支持Map类型
	if p.isTextType {
		return make(map[string]interface{})
	}

	// 获取值并转换，失败返回空Map
	value := p.getValueByPathOrNil(path)
	m, err := p.toMap(value)
	if err != nil {
		return make(map[string]interface{})
	}
	return m
}

// GetStringMap 获取指定路径的map[string]interface{}类型值（同GetMap）
// path：配置路径（如 "a.b.config"）
// 返回：路径对应的Map值；若路径不存在、转换失败或为Text格式，返回空Map（非nil）
func (p *ConfigParser) GetStringMap(path string) map[string]interface{} {
	return p.GetMap(path)
}

// GetStringMapString 获取指定路径的map[string]string类型值
// path：配置路径（如 "a.b.headers"）
// 返回：路径对应的String Map值；若路径不存在、转换失败或为Text格式，返回空Map（非nil）
func (p *ConfigParser) GetStringMapString(path string) map[string]string {
	// Text格式不支持Map类型
	if p.isTextType {
		return make(map[string]string)
	}

	// 先获取Map，再转换为string值的Map
	valueMap := p.GetMap(path)
	result := make(map[string]string, len(valueMap))
	for k, v := range valueMap {
		result[k] = p.convertToString(v)
	}
	return result
}

// GetSlice 获取指定路径的[]interface{}类型值（切片）
// path：配置路径（如 "a.b.list"）
// 返回：路径对应的切片值；若路径不存在、转换失败或为Text格式，返回空切片（非nil）
func (p *ConfigParser) GetSlice(path string) []interface{} {
	// Text格式不支持切片类型
	if p.isTextType {
		return make([]interface{}, 0)
	}

	// 获取值并转换，失败返回空切片
	value := p.getValueByPathOrNil(path)
	slice, err := p.toSlice(value)
	if err != nil {
		return make([]interface{}, 0)
	}
	return slice
}

// GetStringSlice 获取指定路径的[]string类型值（字符串切片）
// path：配置路径（如 "a.b.names"）
// 返回：路径对应的字符串切片；若路径不存在、转换失败或为Text格式，返回空切片（非nil）
func (p *ConfigParser) GetStringSlice(path string) []string {
	// Text格式不支持切片类型
	if p.isTextType {
		return make([]string, 0)
	}

	// 先获取切片，再转换为string元素的切片
	slice := p.GetSlice(path)
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = p.convertToString(v)
	}
	return result
}

// GetIntSlice 获取指定路径的[]int类型值（整数切片）
// path：配置路径（如 "a.b.ids"）
// 返回：路径对应的整数切片；若路径不存在、转换失败或为Text格式，返回空切片（非nil）
func (p *ConfigParser) GetIntSlice(path string) []int {
	// Text格式不支持切片类型
	if p.isTextType {
		return make([]int, 0)
	}

	// 先获取切片，再转换为int元素的切片（单个元素转换失败则整体返回空切片）
	slice := p.GetSlice(path)
	result := make([]int, 0, len(slice))
	for _, v := range slice {
		val, err := p.toInt64(v)
		if err != nil {
			return make([]int, 0)
		}
		result = append(result, int(val))
	}
	return result
}

// ------------------------------
// 私有方法：配置解析（格式-specific）
// ------------------------------

// parseJSON 解析JSON格式配置
func (p *ConfigParser) parseJSON() error {
	return json.Unmarshal([]byte(p.content), &p.data)
}

// parseYAML 解析YAML格式配置（支持yaml/yml后缀）
func (p *ConfigParser) parseYAML() error {
	return yaml.Unmarshal([]byte(p.content), &p.data)
}

// parseTOML 解析TOML格式配置
func (p *ConfigParser) parseTOML() error {
	_, err := toml.Decode(p.content, &p.data)
	return err
}

// parseINI 解析INI格式配置（Section+Key形式，展平为 "Section.Key" 路径）
func (p *ConfigParser) parseINI() error {
	cfg, err := ini.Load([]byte(p.content))
	if err != nil {
		return err
	}

	// 遍历所有Section和Key，展平为 "Section.Key" 格式
	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		for _, key := range section.Keys() {
			fullKey := key.Name()
			// 非默认Section（DEFAULT）需添加Section前缀
			if sectionName != ini.DefaultSection {
				fullKey = fmt.Sprintf("%s.%s", sectionName, fullKey)
			}
			p.data[fullKey] = key.Value()
		}
	}
	return nil
}

// parseProperties 解析Properties格式配置（Key=Value形式，直接使用Key作为路径）
func (p *ConfigParser) parseProperties() error {
	props, err := properties.LoadString(p.content)
	if err != nil {
		return err
	}

	// 遍历所有Key-Value，直接存入data
	for _, key := range props.Keys() {
		value, _ := props.Get(key) // 忽略GetValue错误（默认空串）
		p.data[key] = value
	}
	return nil
}

// ------------------------------
// 私有方法：数据处理与访问
// ------------------------------

// getValueByPathOrNil 通过路径获取配置值，失败返回nil（内部使用，无锁保护）
func (p *ConfigParser) getValueByPathOrNil(path string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// 先查缓存
	if p.cacheEnabled {
		if cached, found := p.cache.Load(path); found {
			return cached
		}
	}

	// 尝试直接匹配展平后的路径（如 "a.b.c"）
	if val, exists := p.data[path]; exists && val != nil {
		if p.cacheEnabled {
			p.cache.Store(path, val)
		}
		return val
	}

	// 拆分路径迭代查找（支持嵌套和数组）
	value, err := p.findNestedValue(p.data, strings.Split(path, "."))
	if err != nil || value == nil {
		return nil
	}

	// 缓存结果
	if p.cacheEnabled {
		p.cache.Store(path, value)
	}
	return value
}

// findNestedValue 迭代查找嵌套路径的值（支持数组），失败返回错误（内部使用，无锁保护）
// data：待查找的数据源（map或切片）
// pathParts：路径拆分后的切片（如 "a.b[0].c" 拆分为 ["a", "b[0]", "c"]）
func (p *ConfigParser) findNestedValue(data interface{}, pathParts []string) (interface{}, error) {
	current := data
	var currentPath []string // 记录当前已处理的路径（用于错误提示）

	for _, part := range pathParts {
		currentPath = append(currentPath, part)
		currentFullPath := p.joinPath(currentPath) // 当前完整路径（如 "a.b[0]"）

		// 处理数组路径（如 "b[0]"）
		if arrayMatch := p.arrayRegex.FindStringSubmatch(part); len(arrayMatch) == 3 {
			arrName := arrayMatch[1] // 数组键名（如 "b"）
			idxStr := arrayMatch[2]  // 数组索引字符串（如 "0"）
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				return nil, fmt.Errorf("路径%s：无效数组索引%s", currentFullPath, idxStr)
			}

			// 检查当前节点是否为Map（需先通过键名获取数组）
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("路径%s：非Map类型，无法访问数组%s", currentFullPath, arrName)
			}

			// 获取数组并验证
			arrVal, exists := currentMap[arrName]
			if !exists {
				return nil, fmt.Errorf("路径%s：数组%s不存在", currentFullPath, arrName)
			}
			arr, ok := arrVal.([]interface{})
			if !ok {
				return nil, fmt.Errorf("路径%s：%s非数组类型", currentFullPath, arrName)
			}
			if idx < 0 || idx >= len(arr) {
				return nil, fmt.Errorf("路径%s：数组%s索引%d越界（长度%d）", currentFullPath, arrName, idx, len(arr))
			}

			// 进入数组的指定索引节点
			current = arr[idx]
			continue
		}

		// 处理普通键（非数组）
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("路径%s：非Map类型，无法访问键%s", currentFullPath, part)
		}
		val, exists := currentMap[part]
		if !exists {
			return nil, fmt.Errorf("路径%s：键%s不存在", currentFullPath, part)
		}

		// 进入下一层节点
		current = val
	}

	return current, nil
}

// flattenData 将嵌套数据展平为路径-值形式（如嵌套map["a"]["b"] -> 路径 "a.b"）
func (p *ConfigParser) flattenData(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	p.flattenRecursive("", data, result)
	return result
}

// flattenRecursive 递归展平数据（内部使用）
// prefix：当前路径前缀（如 "a.b"）
// data：待展平的数据（map/切片/基础类型）
// result：展平后的结果存储（map[路径]值）
func (p *ConfigParser) flattenRecursive(prefix string, data interface{}, result map[string]interface{}) {
	switch v := data.(type) {
	// 处理string键的Map（如JSON/YAML解析结果）
	case map[string]interface{}:
		// 保留原Map的顶层路径（用于GetMap获取完整Map）
		if prefix != "" {
			result[prefix] = v
		}
		// 递归展平子节点
		for key, val := range v {
			newPrefix := p.joinPrefix(prefix, key)
			p.flattenRecursive(newPrefix, val, result)
		}

	// 处理interface键的Map（如YAML解析的特殊情况）
	case map[interface{}]interface{}:
		// 转换为string键的Map并保留顶层路径
		strMap := make(map[string]interface{}, len(v))
		for key, val := range v {
			strKey := fmt.Sprintf("%v", key) // 键转为字符串（如int键1 -> "1"）
			strMap[strKey] = val
		}
		if prefix != "" {
			result[prefix] = strMap
		}
		// 递归展平子节点
		for key, val := range strMap {
			newPrefix := p.joinPrefix(prefix, key)
			p.flattenRecursive(newPrefix, val, result)
		}

	// 处理切片（数组）
	case []interface{}:
		// 保留原切片的顶层路径（用于GetSlice获取完整切片）
		if prefix != "" {
			result[prefix] = v
		}
		// 递归展平切片元素（路径格式："prefix[索引]"）
		for idx, item := range v {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, idx)
			p.flattenRecursive(newPrefix, item, result)
		}

	// 处理基础类型（直接存储）
	case time.Time, time.Duration, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		if prefix != "" {
			result[prefix] = v
		}

	// 其他类型（如nil、struct等，默认转为字符串存储）
	default:
		if prefix != "" {
			result[prefix] = v
		}
	}
}

// ------------------------------
// 私有方法：类型转换（无错误返回时返回默认值）
// ------------------------------

// toInt64 将值转换为int64类型，失败返回0和错误
func (p *ConfigParser) toInt64(value interface{}) (int64, error) {
	if value == nil {
		return 0, fmt.Errorf("值为nil，无法转换为int64")
	}

	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > 1<<63-1 { // 检查是否超出int64最大值
			return 0, fmt.Errorf("uint64值%d超出int64范围", v)
		}
		return int64(v), nil
	case float32:
		if float32(int64(v)) != v { // 检查精度丢失
			return 0, fmt.Errorf("float32值%.2f转换为int64会丢失精度", v)
		}
		return int64(v), nil
	case float64:
		if float64(int64(v)) != v { // 检查精度丢失
			return 0, fmt.Errorf("float64值%.2f转换为int64会丢失精度", v)
		}
		return int64(v), nil
	case string:
		// 支持多进制字符串（0x十六进制、0b二进制、0八进制）
		switch {
		case strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X"):
			return strconv.ParseInt(v[2:], 16, 64)
		case strings.HasPrefix(v, "0b") || strings.HasPrefix(v, "0B"):
			return strconv.ParseInt(v[2:], 2, 64)
		case strings.HasPrefix(v, "0o") || strings.HasPrefix(v, "0O"):
			return strconv.ParseInt(v[2:], 8, 64)
		case strings.HasPrefix(v, "0") && len(v) > 1:
			return strconv.ParseInt(v[1:], 8, 64)
		default:
			return strconv.ParseInt(v, 10, 64)
		}
	default:
		return 0, fmt.Errorf("类型%T无法转换为int64", value)
	}
}

// toUint64 将值转换为uint64类型，失败返回0和错误
func (p *ConfigParser) toUint64(value interface{}) (uint64, error) {
	if value == nil {
		return 0, fmt.Errorf("值为nil，无法转换为uint64")
	}

	switch v := value.(type) {
	case int:
		if v < 0 {
			return 0, fmt.Errorf("负值%d无法转换为uint64", v)
		}
		return uint64(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("负值%d无法转换为uint64", v)
		}
		return uint64(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("负值%d无法转换为uint64", v)
		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("负值%d无法转换为uint64", v)
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("负值%d无法转换为uint64", v)
		}
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case float32:
		if v < 0 {
			return 0, fmt.Errorf("负值%.2f无法转换为uint64", v)
		}
		if float32(uint64(v)) != v { // 检查精度丢失
			return 0, fmt.Errorf("float32值%.2f转换为uint64会丢失精度", v)
		}
		return uint64(v), nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("负值%.2f无法转换为uint64", v)
		}
		if float64(uint64(v)) != v { // 检查精度丢失
			return 0, fmt.Errorf("float64值%.2f转换为uint64会丢失精度", v)
		}
		return uint64(v), nil
	case string:
		// 支持多进制字符串（0x十六进制、0b二进制、0八进制）
		switch {
		case strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X"):
			return strconv.ParseUint(v[2:], 16, 64)
		case strings.HasPrefix(v, "0b") || strings.HasPrefix(v, "0B"):
			return strconv.ParseUint(v[2:], 2, 64)
		case strings.HasPrefix(v, "0o") || strings.HasPrefix(v, "0O"):
			return strconv.ParseUint(v[2:], 8, 64)
		case strings.HasPrefix(v, "0") && len(v) > 1:
			return strconv.ParseUint(v[1:], 8, 64)
		default:
			return strconv.ParseUint(v, 10, 64)
		}
	default:
		return 0, fmt.Errorf("类型%T无法转换为uint64", value)
	}
}

// toFloat64 将值转换为float64类型，失败返回0.0和错误
func (p *ConfigParser) toFloat64(value interface{}) (float64, error) {
	if value == nil {
		return 0.0, fmt.Errorf("值为nil，无法转换为float64")
	}

	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0.0, fmt.Errorf("类型%T无法转换为float64", value)
	}
}

// toBool 将值转换为bool类型，失败返回false和错误
func (p *ConfigParser) toBool(value interface{}) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("值为nil，无法转换为bool")
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil // 非0即true
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("类型%T无法转换为bool", value)
	}
}

// toTime 将值转换为time.Time类型，失败返回时间零值和错误
func (p *ConfigParser) toTime(value interface{}) (time.Time, error) {
	if value == nil {
		return time.Time{}, fmt.Errorf("值为nil，无法转换为time.Time")
	}

	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		// 支持常见时间格式（按优先级匹配）
		formats := []string{
			time.RFC3339,          // 2006-01-02T15:04:05Z07:00
			time.RFC3339Nano,      // 2006-01-02T15:04:05.999999999Z07:00
			"2006-01-02 15:04:05", // 年月日 时分秒
			"2006-01-02",          // 年月日
			"15:04:05",            // 时分秒
			time.RFC1123,          // Mon, 02 Jan 2006 15:04:05 MST
			time.RFC1123Z,         // Mon, 02 Jan 2006 15:04:05 -0700
			time.RFC822,           // 02 Jan 06 15:04 MST
			time.RFC822Z,          // 02 Jan 06 15:04 -0700
		}

		var lastErr error
		for _, format := range formats {
			t, err := time.Parse(format, v)
			if err == nil {
				return t, nil
			}
			lastErr = err
		}
		return time.Time{}, fmt.Errorf("字符串\"%s\"无法解析为时间: %w", v, lastErr)
	default:
		return time.Time{}, fmt.Errorf("类型%T无法转换为time.Time", value)
	}
}

// toDuration 将值转换为time.Duration类型，失败返回0和错误
func (p *ConfigParser) toDuration(value interface{}) (time.Duration, error) {
	if value == nil {
		return 0, fmt.Errorf("值为nil，无法转换为time.Duration")
	}

	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		return time.ParseDuration(v)
	default:
		return 0, fmt.Errorf("类型%T无法转换为time.Duration", value)
	}
}

// toMap 将值转换为map[string]interface{}类型，失败返回nil和错误
func (p *ConfigParser) toMap(value interface{}) (map[string]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("值为nil，无法转换为map[string]interface{}")
	}

	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	case map[interface{}]interface{}:
		// 转换interface键为string（如YAML解析的Map）
		result := make(map[string]interface{}, len(v))
		for k, val := range v {
			result[fmt.Sprintf("%v", k)] = val
		}
		return result, nil
	default:
		return nil, fmt.Errorf("类型%T无法转换为map[string]interface{}", value)
	}
}

// toSlice 将值转换为[]interface{}类型，失败返回nil和错误
func (p *ConfigParser) toSlice(value interface{}) ([]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("值为nil，无法转换为[]interface{}")
	}

	slice, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("类型%T无法转换为[]interface{}", value)
	}
	return slice, nil
}

// convertToString 将任意类型值转换为字符串（无失败，默认返回fmt.Sprintf("%v", value)）
func (p *ConfigParser) convertToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32) // 保留原始精度
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339) // 统一时间字符串格式
	case time.Duration:
		return v.String()
	case []interface{}:
		// 切片转为JSON数组字符串（失败则返回默认格式）
		if jsonBytes, err := json.Marshal(v); err == nil {
			return string(jsonBytes)
		} else {
			log.Printf("警告：切片序列化JSON失败（值：%v）：%v", v, err)
		}
	case map[string]interface{}, map[interface{}]interface{}:
		// Map转为JSON对象字符串（失败则返回默认格式）
		if jsonBytes, err := json.Marshal(v); err == nil {
			return string(jsonBytes)
		} else {
			log.Printf("警告：Map序列化JSON失败（值：%v）：%v", v, err)
		}
	}

	// 其他类型：默认转为字符串（如struct、nil等）
	return fmt.Sprintf("%v", value)
}

// joinPrefix 拼接路径前缀和键（如 prefix="a.b"，key="c" → "a.b.c"）
func (p *ConfigParser) joinPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", prefix, key)
}

// joinPath 拼接路径片段为完整路径（如 ["a", "b[0]", "c"] → "a.b[0].c"）
func (p *ConfigParser) joinPath(parts []string) string {
	return strings.Join(parts, ".")
}
