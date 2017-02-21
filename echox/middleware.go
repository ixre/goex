/**
 * Copyright 2015 @ at3.net.
 * name : middleware.go
 * author : jarryliu
 * date : 2016-11-20 22:42
 * description :
 * history :
 */
package echox

import (
	"github.com/labstack/echo"
	"net/http"
	"regexp"
)

var (
	requestFilter = map[string]*regexp.Regexp{
		"GET": regexp.MustCompile("'|(and\\s|or\\s)\\b.+?(>|<|=|in|like)|\\/\\*" +
			".+?\\*\\/|<\\s*script\\b|\\bEXEC\\b|UNION" +
			".+?SELECT|UPDATE.+?SET|INSERT\\s+INTO.+?VALUES|(SELECT|DELETE).+?FROM|(CREATE|ALTER|DROP|TRUNCATE)\\s+" +
			"(TABLE|DATABASE)"),
		"POST": regexp.MustCompile("\\b(and\\s|or\\s)\\b.{1,6}?(=|>|<|\\bin\\b|\\blike\\b)|\\/\\*" +
			".+?\\*\\/|<\\s*script\\b|\\bEXEC\\b|UNION.+?SELECT|UPDATE.+?SET|INSERT\\s+INTO.+?VALUES|(SELECT|DELETE).+?FROM|" +
			"(CREATE|ALTER|DROP|TRUNCATE)\\s+(TABLE|DATABASE)"),
	}

	/*
	   getFilter = postFilter = cookieFilter = regexp.MustCompile("\\b(and|or)\\b.{1,6}?(=|>|<|\\bin\\b|\\blike\\b)|\\/\\*.+?\\*\\/|<\\s*script\\b|\\bEXEC\\b|UNION.+?SELECT|UPDATE.+?SET|INSERT\\s+INTO.+?VALUES|(SELECT|DELETE).+?FROM|(CREATE|ALTER|DROP|TRUNCATE)\\s+(TABLE|DATABASE)");
	*/
)

// 防SQL注入
func StopAttackMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		badRequest := false
		method := req.Method
		switch method {
		case "GET":
			badRequest = requestFilter[method].MatchString(req.URL.RawQuery)
		case "POST":
			badRequest = requestFilter["GET"].MatchString(req.URL.RawQuery) ||
				requestFilter[method].MatchString(
					req.Form.Encode())
		}
		if badRequest {
			return c.HTML(http.StatusNotFound,
				"<div style='color:red;'>您提交的参数非法,系统已记录您本次操作!</div>")
		}
		return h(c)
	}
}
