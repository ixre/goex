package generator

import (
	"github.com/jsix/gof/util"
	ht "html/template"
	"strings"
	"unicode"
)

type internalFunc struct {
}

// 返回模板函数
func (t *internalFunc) funcMap() ht.FuncMap {
	fm := make(map[string]interface{})
	fm["boolInt"] = t.boolInt
	fm["isEmpty"] = t.isEmpty
	fm["rawHtml"] = t.rawHtml
	fm["add"] = t.plus
	fm["plus"] = t.plus
	fm["multi"] = t.multi
	fm["mathRemain"] = t.mathRemain
	fm["lowerTitle"] = t.lowerTitle
	fm["title"] = t.title
	fm["lower"] = t.lower
	fm["upper"] = t.upper
	return fm
}

// 小写
func (t *internalFunc) lower(s string) string {
	return strings.ToLower(s)
}

// 大写
func (t *internalFunc) upper(s string) string {
	return strings.ToUpper(s)
}

// 将首字母小写
func (t *internalFunc) lowerTitle(s string) string {
	if rune0 := rune(s[0]); unicode.IsUpper(rune0) {
		return string(unicode.ToLower(rune0)) + s[1:]
	}
	return s
}

// 将字符串单词首字母大写
func (t *internalFunc) title(s string) string {
	return strings.Title(s)
}

// 判断是否为true
func (t *internalFunc) boolInt(i int32) bool {
	return i > 0
}

// 加法
func (t *internalFunc) plus(x, y int) int {
	return x + y
}

// 乘法
func (t *internalFunc) multi(x, y interface{}) interface{} {
	fx, ok := x.(float64)
	if ok {
		switch y.(type) {
		case float32:
			return fx * float64(y.(float32))
		case float64:
			return fx * y.(float64)
		case int:
			return fx * float64(y.(int))
		case int32:
			return fx * float64(y.(int32))
		case int64:
			return fx * float64(y.(int64))
		}
	}
	panic("not support")
}

// I32转为字符
func (t *internalFunc) str(i interface{}) string {
	return util.Str(i)
}

// 是否为空
func (t *internalFunc) isEmpty(s string) bool {
	if s == "" {
		return true
	}
	return strings.TrimSpace(s) == ""
}

// 转换为HTML
func (t *internalFunc) rawHtml(v interface{}) ht.HTML {
	return ht.HTML(util.Str(v))
}

//求余
func (t *internalFunc) mathRemain(i int, j int) int {
	return i % j
}
