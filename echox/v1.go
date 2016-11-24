/**
 * Copyright 2015 @ at3.net.
 * name : v1.go
 * author : jarryliu
 * date : 2016-11-20 22:17
 * description :
 * history :
 */
package echox

import (
	"net/http"
	"reflect"
	"strings"
	"sync"
)

var (
	mux             sync.Mutex
	dynamicHandlers = make(map[string]Handler) // 动态处理程序
)

// get handler by reflect
func getHandler(v interface{}, action string) (Handler, bool) {
	t := reflect.ValueOf(v)
	method := t.MethodByName(strings.Title(action))
	if method.IsValid() {
		v, ok := method.Interface().(func(*Context) error)
		return v, ok
	}
	return nil, false
}

func getMvcHandler(route string, c *Context, obj interface{}) Handler {
	mux.Lock()
	defer mux.Unlock()
	a := c.Param("action")
	k := route + a
	if dynamicHandlers == nil {
		dynamicHandlers = make(map[string]Handler)
	}
	if v, ok := dynamicHandlers[k]; ok {
		//查找路由表
		return v
	}
	if v, ok := getHandler(obj, a); ok {
		//存储路由表
		dynamicHandlers[k] = v
		return v
	}
	return nil
}

// 注册动态获取处理程序, 请使用e.Auto(path,handler)代替
func (e *Echo) XaAny(path string, obj interface{}) {
	h := func(c *Context) error {
		if c.Param("action") == "" {
			return c.String(http.StatusInternalServerError, "route must contain :action")
		}
		if hd := getMvcHandler(path, c, obj); hd != nil {
			return hd(c)
		}
		return c.String(http.StatusNotFound, "no such file")
	}
	e.Any(path, h)
}
