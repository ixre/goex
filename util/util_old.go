/**
 * Copyright 2015 @ z3q.net.
 * name : string.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"encoding/json"
	"html/template"
	"log"
)

// 强制序列化为可用于HTML的JSON
func MustHtmlJson(v interface{}) template.JS {
	d, err := json.Marshal(v)
	if err != nil {
		log.Println("[ Json][ Mashal]: ", err.Error())
	}
	return template.JS(d)
}
