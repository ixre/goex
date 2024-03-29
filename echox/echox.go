/**
 * Copyright 2015 @ 56x.net.
 * name : echo
 * author : jarryliu
 * date : 2015-12-04 10:51
 * description :
 * history :
 */
package echox

import (
	"errors"
	"github.com/ixre/gof"
	"github.com/ixre/gof/log"
	http2 "github.com/ixre/gof/net/http"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/web"
	"github.com/ixre/gof/web/session"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"reflect"
	"strings"
	"unicode"
)

var (
	_ echo.Renderer = new(renderer)
	_ echo.Context  = new(Context)
)

type (
	Echo struct {
		*echo.Echo
		app    gof.App
		varMap map[string]interface{}
	}
	Group struct {
		*echo.Group
		echo *Echo
	}
	Context struct {
		echo.Context
		echo    *Echo
		App     gof.App
		Session *session.Session
		Storage storage.Interface
	}
	TemplateData struct {
		Var  map[string]interface{}
		Map  map[string]interface{}
		Data interface{}
	}
	Handler func(*Context) error
)

type renderer struct {
	t *gof.CacheTemplate
}

func NewRenderer(t *gof.CacheTemplate) echo.Renderer {
	return &renderer{t}
}

func (g *renderer) Render(w io.Writer, name string,
	data interface{}, c echo.Context) error {
	return g.t.Execute(w, name, data)
}

// new echo instance
func New() *Echo {
	e := &Echo{
		Echo:   echo.New(),
		varMap: make(map[string]interface{}),
	}
	e.Echo.Use(e.contextMiddle)
	e.Echo.HTTPErrorHandler = e.defaultErrorHandler
	return e
}

// 默认错误处理程序
func (e *Echo) defaultErrorHandler(err error, c echo.Context) {
	msg := err.Error()
	i := strings.Index(msg, "message=")
	if i != -1 {
		msg = msg[i+8:]
	}
	c.String(http.StatusInternalServerError, msg)
	e.Echo.Logger.Error(err)
}

// 上下文中间处理,兼容*Context和echo.Context作为函数签名
func (e *Echo) contextMiddle(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, ok := c.(*Context); ok {
			return h(c)
		}
		return h(e.parseContext(c, e.app))
	}
}

// 转换为Echo Handler
func (e *Echo) parseHandler(h Handler) func(c echo.Context) error {
	if e.app == nil {
		if gof.CurrentApp == nil {
			panic(errors.New("not register or no global serve instance for echox!"))
		}
		e.app = gof.CurrentApp
	}
	return func(c echo.Context) error {
		return h(e.parseContext(c, e.app))
	}
}

func (e *Echo) parseContext(c echo.Context, app gof.App) *Context {
	req, rsp := c.Request(), c.Response()
	s := session.Default(rsp, req)
	return &Context{
		Context: c,
		echo:    e,
		Session: s,
		Storage: app.Storage(),
		App:     app,
	}
}

func (e *Echo) SetApp(app gof.App) {
	e.app = app
}

// 获取原始的Echo对象
func (e *Echo) Original() *echo.Echo {
	return e.Echo
}

// 分组
func (e *Echo) Group(prefix string, m ...echo.MiddlewareFunc) *Group {
	g := e.Echo.Group(prefix, m...)
	return &Group{
		Group: g,
		echo:  e,
	}
}

// 设置模板
func (e *Echo) SetRenderer(r echo.Renderer) {
	e.Renderer = r
}

// 设置变量
func (e *Echo) SetVariable(key string, v interface{}) {
	e.varMap[key] = v
}

// 获取变量
func (e *Echo) GetVariable(key string) interface{} {
	return e.varMap[key]
}

// 获取Echo原始对象
func (e *Echo) Classic() *echo.Echo {
	return e.Echo
}

// 静态文件路由
func (e *Echo) Static(prefix, root string) {
	// 解决prefix以"/"结尾，无法工作的BUG
	if l := len(prefix); prefix[l-1] == '/' {
		prefix = prefix[:l-1]
	}
	e.Echo.Static(prefix, root)
}

// 注册自定义的GET处理程序
func (e *Echo) GET(path string, h Handler) {
	e.Echo.GET(path, e.parseHandler(h))
}

// 注册自定义的POST处理程序
func (e *Echo) POST(path string, h Handler) {
	e.Echo.POST(path, e.parseHandler(h))
}

// 注册自定义的GET/POST处理程序
func (e *Echo) Any(path string, h Handler) {
	e.Echo.Any(path, e.parseHandler(h))
}

// 将控制器下所有的动作映射到路由,如果动作名只有首字母为大写，
// 那么URL中动作名小写，反之与动作名一致
func (e *Echo) Auto(prefix string, i interface{}) {
	mp := getHandlerArray(i)
	for k, v := range mp {
		e.Any(prefix+"/"+k, v)
	}
}

func (e *Echo) AutoGET(prefix string, i interface{}) {
	mp := getHandlerArray(i)
	for k, v := range mp {
		e.GET(prefix+"/"+k, v)
	}
}

func (e *Echo) AutoPOST(prefix string, i interface{}) {
	mp := getHandlerArray(i)
	for k, v := range mp {
		e.POST(prefix+"/"+k, v)
	}
}

// 获取Echo原始的Group对象
func (g *Group) Classic() *echo.Group {
	return g.Group
}

func (g *Group) GET(path string, h Handler) {
	g.Group.GET(path, g.echo.parseHandler(h))
}

func (g *Group) POST(path string, h Handler) {
	g.Group.POST(path, g.echo.parseHandler(h))
}

func (g *Group) Any(path string, h Handler) {
	g.Group.Any(path, g.echo.parseHandler(h))
}

// 将控制器下所有的动作映射到路由
func (g *Group) Auto(prefix string, i interface{}) {
	mp := getHandlerArray(i)
	for k, v := range mp {
		g.Group.Any(prefix+"/"+k, g.echo.parseHandler(v))
	}
}
func (g *Group) AutoGET(prefix string, i interface{}) {
	mp := getHandlerArray(i)
	for k, v := range mp {
		g.Group.GET(prefix+"/"+k, g.echo.parseHandler(v))
	}
}

func (g *Group) AutoPOST(prefix string, i interface{}) {
	mp := getHandlerArray(i)
	panic(len(mp))
	for k, v := range mp {
		g.Group.POST(prefix+"/"+k, g.echo.parseHandler(v))
	}
}

/*********  以下需重构   **********/

// 返回真实IP
func (c *Context) RealIp() string {
	return http2.RealIp(c.Request())
}
func (c *Context) IsPost() bool {
	return c.Request().Method == "POST"
}

// 获取请求完整的地址
func (c *Context) RequestRawURI(r *http.Request) string {
	return web.RequestRawURI(r)
}

func (c *Context) StringOK(s string) error {
	return c.debug(c.String(http.StatusOK, s))
}

func (c *Context) Error(err error) {
	if err != nil {
		//web.HttpError(c.Response(), err)
		c.Response().Write([]byte(err.Error()))
	}
}

func (c *Context) debug(err error) error {
	if err != nil {
		web.HttpError(c.Response(), err)
		return nil
	}
	return err
}

func (c *Context) Debug(err error) error {
	return c.debug(err)
}

// 覆写Render方法
func (c *Context) Render(code int, name string, data interface{}) error {
	return c.debug(c.Context.Render(code, name, data))
}

func (c *Context) RenderOK(name string, data interface{}) error {
	return c.debug(c.Render(http.StatusOK, name, data))
}

func (c *Context) NewData() *TemplateData {
	return &TemplateData{
		Var:  c.echo.varMap,
		Map:  make(map[string]interface{}),
		Data: nil,
	}
}

// 获取控制器所有的动作映射
func getHandlerArray(i interface{}) map[string]Handler {
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		log.Println("[ echox][ warning]: ", t.String(), " not a pointer.")
	}
	mp := map[string]Handler{}
	for k, j := 0, v.NumMethod(); k < j; k++ {
		m := v.Method(k)
		if !m.IsValid() {
			continue
		}
		v2, ok := m.Interface().(func(*Context) error)
		if ok {
			name := unRouteTitle(t.Method(k).Name)
			mp[name] = v2
		}
	}
	return mp
}

//如果除首字母外均为为小写，则小写
func unRouteTitle(s string) string {
	first := unicode.ToLower(rune(s[0]))
	r := append([]rune{first}, []rune(s[1:])...)
	return string(r)
}
