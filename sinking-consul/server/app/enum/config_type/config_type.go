package config_type

const (
	TEXT       = "text"
	JSON       = "json"
	YAML       = "yaml"
	INI        = "ini"
	TOML       = "toml"
	PROPERTIES = "properties"
)

// Map 类型数据
func Map() map[string]string {
	return map[string]string{
		TEXT:       "TEXT",
		JSON:       "JSON",
		YAML:       "YAML",
		INI:        "INI",
		TOML:       "TOML",
		PROPERTIES: "Properties",
	}
}
