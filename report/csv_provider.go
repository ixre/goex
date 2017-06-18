package report

import (
	"bytes"
	"log"
	"strings"
)

var _ IDataExportProvider = new(CsvProvider)

type CsvProvider struct {
	delimer string
}

func NewCsvProvider() IDataExportProvider {
	return &CsvProvider{
		delimer: " ",
	}
}

func (c *CsvProvider) Export(rows []map[string]interface{},
	keys []string, alias []string) (binary []byte) {
	buf := bytes.NewBufferString("")
	// 显示表头
	showHeader := keys != nil && len(keys) > 0
	log.Println("---", showHeader)
	if showHeader {
		for i, k := range alias {
			if i > 0 {
				buf.WriteString(c.delimer)
			}
			buf.WriteString(k)
		}
	}
	l := len(rows)
	for i, row := range rows {
		if i < l {
			buf.WriteString("\n")
		}
		for ki, k := range keys {
			if ki > 0 {
				buf.WriteString(c.delimer)
			}
			data := row[k].(string)
			//防止里面含有特殊符号
			data = strings.Replace(data, "\"", "\"\"", -1)
			buf.WriteString("\"")
			buf.WriteString(data)
			buf.WriteString("\"")
		}
	}
	return buf.Bytes()
}
