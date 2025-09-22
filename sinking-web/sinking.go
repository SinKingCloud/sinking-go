package sinking_web

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"path"
	"strings"
	"time"
)

type H map[string]interface{}

type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine
	}

	ErrorHandel struct {
		NotFound func(*Context)
		Fail     func(c *Context, code int, message string)
	}

	Engine struct {
		*RouterGroup
		router             *router
		errorHandel        *ErrorHandel
		groups             []*RouterGroup
		htmlTemplates      *template.Template
		funcMap            template.FuncMap
		MaxMultipartMemory int64
		debug              bool
		readTimeout        time.Duration
		writeTimeout       time.Duration
	}
)

func New() *Engine {
	engine := &Engine{
		router:             newRouter(),
		MaxMultipartMemory: defaultMultipartMemory,
		debug:              false,
		readTimeout:        time.Second * 60,
		writeTimeout:       time.Second * 60,
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery())
	if engine.debug {
		engine.Use(Logger())
	}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) *RouterGroup {
	group.middlewares = append(group.middlewares, middlewares...)
	return group
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) ANY(pattern string, handler HandlerFunc) *RouterGroup {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodTrace,
		http.MethodPatch,
		http.MethodHead,
	}
	for _, v := range methods {
		group.addRoute(v, pattern, handler)
	}
	return group
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodGet, pattern, handler)
	return group
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodPost, pattern, handler)
	return group
}

func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodOptions, pattern, handler)
	return group
}

func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodPut, pattern, handler)
	return group
}

func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodDelete, pattern, handler)
	return group
}

func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodHead, pattern, handler)
	return group
}

func (group *RouterGroup) TRACE(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodTrace, pattern, handler)
	return group
}

func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) *RouterGroup {
	group.addRoute(http.MethodPatch, pattern, handler)
	return group
}

func (group *RouterGroup) SetErrorHandle(handle *ErrorHandel) *RouterGroup {
	group.engine.errorHandel = handle
	return group
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) *RouterGroup {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
	return group
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) *Engine {
	engine.funcMap = funcMap
	return engine
}

func (engine *Engine) LoadHtmlGlob(pattern string) *Engine {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
	return engine
}

// SetDebugMode 设置运行模式为debug
func (engine *Engine) SetDebugMode(mode bool) *Engine {
	engine.debug = mode
	return engine
}

// SetTimeOut 设置超时时间
func (engine *Engine) SetTimeOut(read time.Duration, write time.Duration) *Engine {
	engine.readTimeout = read
	engine.writeTimeout = write
	return engine
}

func (group *RouterGroup) PROXY(pattern string, uri string, logger *log.Logger, filter func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy), errorHandle func(http.ResponseWriter, *http.Request, error)) *RouterGroup {
	fun := func(c *Context) {
		prefix := uri[0:2]
		if prefix == "ws" {
			fun := func(r *http.Request, w http.ResponseWriter) {
				if filter != nil {
					filter(r, w, nil)
				}
			}
			_ = c.WebSocketProxy(uri, logger, fun, errorHandle)
		} else {
			_ = c.HttpProxy(uri, logger, filter, errorHandle)
		}
	}
	group.ANY(pattern, fun)
	return group
}

func server(addr string, engine *Engine) *http.Server {
	Author(engine, addr)
	server := &http.Server{
		ReadTimeout:  engine.readTimeout,
		WriteTimeout: engine.writeTimeout,
		Addr:         addr,
		Handler:      engine,
	}
	return server
}

func (engine *Engine) Run(addr string) (err error) {
	server := server(addr, engine)
	return server.ListenAndServe()
}

func (engine *Engine) RunTLS(addr string, certFile string, keyFile string) (err error) {
	server := server(addr, engine)
	return server.ListenAndServeTLS(certFile, keyFile)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	req.URL.Path = engine.cleanPath(req.URL.Path)
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

func (engine *Engine) cleanPath(p string) string {
	const stackBufSize = 128
	if p == "" {
		return "/"
	}
	buf := make([]byte, 0, stackBufSize)
	n := len(p)
	r := 1
	w := 1
	if p[0] != '/' {
		r = 0
		if n+1 > stackBufSize {
			buf = make([]byte, n+1)
		} else {
			buf = buf[:n+1]
		}
		buf[0] = '/'
	}
	trailing := n > 1 && p[n-1] == '/'
	for r < n {
		switch {
		case p[r] == '/':
			r++
		case p[r] == '.' && r+1 == n:
			trailing = true
			r++
		case p[r] == '.' && p[r+1] == '/':
			r += 2
		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			r += 3
			if w > 1 {
				w--
				if len(buf) == 0 {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}
		default:
			if w > 1 {
				engine.bufApp(&buf, p, w, '/')
				w++
			}
			for r < n && p[r] != '/' {
				engine.bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}
	if trailing && w > 1 {
		engine.bufApp(&buf, p, w, '/')
		w++
	}
	if len(buf) == 0 {
		return p[:w]
	}
	return string(buf[:w])
}

func (engine *Engine) bufApp(buf *[]byte, s string, w int, c byte) {
	b := *buf
	if len(b) == 0 {
		if s[w] == c {
			return
		}
		length := len(s)
		if length > cap(b) {
			*buf = make([]byte, length)
		} else {
			*buf = (*buf)[:length]
		}
		b = *buf
		copy(b, s[:w])
	}
	b[w] = c
}
