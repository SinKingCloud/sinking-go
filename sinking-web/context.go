package sinking_web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

// Context 上下文结构体
type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	handlers   []HandlerFunc
	index      int
	engine     *Engine
	lock       sync.RWMutex
	Keys       map[string]interface{}
}

// newContext 实例化一个新的上下文
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:    req.URL.Path,
		Method:  req.Method,
		Request: req,
		Writer:  w,
		index:   -1,
	}
}

// Next 中间件执行
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// IsAborted 是否被挂起
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort 阻止挂起函数
func (c *Context) Abort() {
	c.index = abortIndex
}

// Fail 输出失败响应
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	if c.engine.errorHandel != nil && c.engine.errorHandel.Fail != nil {
		c.engine.errorHandel.Fail(c, code, err)
	} else {
		c.JSON(code, H{"code": code, "message": err})
	}
}

// AllParam 获取所有路径参数
func (c *Context) AllParam() map[string]string {
	return c.Params
}

// Param 获取单个路径参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// DefaultParam 获取路径参数，不存在返回默认参数
func (c *Context) DefaultParam(key, defaultValue string) string {
	value, exists := c.Params[key]
	if exists {
		return value
	}
	return defaultValue
}

// AllForm 获取所有post参数
func (c *Context) AllForm() map[string]string {
	param := map[string]string{}
	err := c.Request.ParseForm()
	if err != nil {
		return param
	}
	for k, v := range c.Request.PostForm {
		param[k] = v[0]
	}
	return param
}

// Form 获取单个post参数
func (c *Context) Form(key string) string {
	return c.Request.FormValue(key)
}

// DefaultForm 获取单个post参数，不存在则返回默认参数
func (c *Context) DefaultForm(key, defaultValue string) string {
	if value := c.Request.FormValue(key); value != "" {
		return value
	}
	return defaultValue
}

// AllQuery 获取所有get参数
func (c *Context) AllQuery() map[string]string {
	param := map[string]string{}
	for k, v := range c.Request.URL.Query() {
		param[k] = v[0]
	}
	return param
}

// Query 获取单个get参数
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// DefaultQuery 获取单个get参数，不存在则返回默认值
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value := c.Request.URL.Query().Get(key); value != "" {
		return value
	}
	return defaultValue
}

// Body 获取请求体
func (c *Context) Body() string {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

// Redirect 重定向跳转
func (c *Context) Redirect(code int, location string) {
	c.SetHeader(HeaderLocation, location)
	c.SetStatus(code)
}

// FormFile 获取上传文件
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	return fh, err
}

// MultipartForm 多文件上传
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(c.engine.MaxMultipartMemory)
	return c.Request.MultipartForm, err
}

// SaveUploadedFile 保存上传文件
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)
	_, err = io.Copy(out, src)
	return err
}

// SetStatus 设置http响应码
func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	if c.StatusCode == 0 {
		c.Writer.WriteHeader(code)
	}
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 返回字符串内容
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader(ContentType, ContentTypeText)
	c.SetStatus(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		return
	}
}

// JSON 返回json内容
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader(ContentType, ContentTypeJson)
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

// Data 返回字符内容
func (c *Context) Data(code int, data []byte) {
	c.SetStatus(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		return
	}
}

// HTML 返回html内容
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader(ContentType, ContentTypeHtml)
	c.SetStatus(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
		return
	}

}

// Set 中间件设置传递内容
func (c *Context) Set(key string, value interface{}) {
	c.lock.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
	c.lock.Unlock()
}

// Get 中间件获取传递内容
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.lock.RLock()
	value, exists = c.Keys[key]
	c.lock.RUnlock()
	return value, exists
}

////  以下为获取不同类型的中间件参数内容  \\\\

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *Context) GetStringMapStringSlice(key string) (temp map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		temp, _ = val.(map[string][]string)
	}
	return
}
