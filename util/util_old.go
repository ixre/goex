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

/*
//编码
func EncodingTransform(src []byte, enc string) ([]byte, error) {
	var ec encoding.Encoding
	switch enc {
	default:
		return src, nil
	case "GBK":
		ec = simplifiedchinese.GBK
	case "GB2312":
		ec = simplifiedchinese.HZGB2312
	case "BIG5":
		ec = traditionalchinese.Big5
	}
	dst := make([]byte, len(src)*2)
	n, _, err := ec.NewEncoder().Transform(dst, src, true)
	return dst[:n], err
}
*/

// 强制序列化为可用于HTML的JSON
func MustHtmlJson(v interface{}) template.JS {
	d, err := json.Marshal(v)
	if err != nil {
		log.Println("[ Json][ Mashal]: ", err.Error())
	}
	return template.JS(d)
}
