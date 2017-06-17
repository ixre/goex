package report

import (
	"bytes"
	"strings"
)

var _ IDataExportProvider = new(TextProvider)

type TextProvider struct {
	delimer string
}

func NewTextProvider() IDataExportProvider {
	return &TextProvider{
		delimer: " ",
	}
}

func (t *TextProvider) Export(rows []map[string]interface{},
	keys []string, alias []string) (binary []byte) {
	buf := bytes.NewBufferString("")
	// 显示表头
	showHeader := keys != nil && len(keys) > 0
	if showHeader {
		for i, k := range alias {
			if i > 0 {
				buf.WriteString(t.delimer)
			}
			buf.WriteString(k)
		}
	}
	l := len(rows)
	for i, row := range rows {
		if i < l-1 {
			buf.WriteString("\n")
		}
		for ki, k := range keys {
			if ki > 0 {
				buf.WriteString(t.delimer)
			}
			data := string(row[k].([]byte))
			//防止里面含有特殊符号
			data = strings.Replace(data, "\"", "\"\"", -1)
			buf.WriteString("\"")
			buf.WriteString(data)
			buf.WriteString("\"")
		}
	}
	return buf.Bytes()
}
