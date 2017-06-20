package report

import (
	"bytes"
	"strings"
)

var _ IDataExportProvider = new(CsvProvider)

type CsvProvider struct {
	delimer string
}

func NewCsvProvider() IDataExportProvider {
	return &CsvProvider{
		delimer: ",",
	}
}

func (c *CsvProvider) Export(rows []map[string]interface{},
	fields []string, names []string, formatter IExportFormatter) (binary []byte) {
	buf := bytes.NewBufferString("")
	// 显示表头
	showHeader := fields != nil && len(fields) > 0
	if showHeader {
		for i, k := range names {
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
		for ki, k := range fields {
			if ki > 0 {
				buf.WriteString(c.delimer)
			}
			if row[k] == nil {
				buf.WriteString("")
				continue
			}
			if formatter != nil {
				row[k] = formatter.Format(k, names[ki], row[k])
			}
			data := row[k].(string)
			specData := strings.Index(data, " ") != -1 ||
				strings.Index(data, "-") != -1 ||
				strings.Index(data, "'") != -1

			if strings.Index(data, "\"") != -1 {
				data = strings.Replace(data, "\"", "\"\"", -1)
				specData = true
			}
			//防止里面含有特殊符号
			if specData {
				buf.WriteString("\"")
				buf.WriteString(data)
				buf.WriteString("\"")
			} else {
				buf.WriteString(data)
			}
		}
	}
	return buf.Bytes()
}
