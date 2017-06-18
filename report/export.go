/**
 * Copyright 2013 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package report

import (
	"database/sql"
	_ "database/sql"
	"encoding/xml"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
)

type (
	IDbProvider interface {
		//获取数据库连接
		GetDB() *sql.DB
	}

	//数据项
	DataExportPortal struct {
	}

	//列映射
	ColumnMapping struct {
		//列的字段
		Field string
		//列的名称
		Name string
	}

	//导入导出项目配置
	ItemConfig struct {
		ColumnMapping string
		Query         string
		Total         string
		Import        string
	}

	//数据导出入口
	IDataExportPortal interface {
		//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
		GetColumnNames() []ColumnMapping

		//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
		//ColumnNames() (names []DataColumnMapping)
		//获取要导出的数据及表结构
		GetSchemaAndData(ht map[string]string) (rows []map[string]interface{}, total int, err error)
		//获取要导出的数据Json格式
		GetJsonData(ht map[string]string) string
		//获取统计数据
		GetTotalView(ht map[string]string) (row map[string]interface{})
		//根据导出的列名获取列的索引及对应键
		GetExportColumnNames(exportColumnNames []string) (fields []string)
	}

	//导出
	IDataExportProvider interface {
		//导出
		Export(rows []map[string]interface{}, keys []string, alias []string) (binary []byte)
	}

	//导出参数
	Params struct {
		//参数
		Params map[string]string
		//要到导出的列的名称集合
		ExportColumnNames []string
	}
)

// 从Map中拷贝数据
func (p *Params) Copy(form map[string]string) {
	for k, v := range form {
		if k != "total" && k != "rows" && k != "params" {
			p.Params[k] = strings.TrimSpace(v)
		}
	}
}

// 从表单参数中导入数据
func (p *Params) CopyForm(form url.Values) {
	for k, v := range form {
		if k != "total" && k != "rows" && k != "params" {
			p.Params[k] = strings.TrimSpace(v[0])
		}
	}
}

//获取列映射数组
func readItemConfigFromXml(xmlFilePath string) (*ItemConfig, error) {
	var cfg ItemConfig
	content, _err := ioutil.ReadFile(xmlFilePath)
	if _err != nil {
		return &ItemConfig{}, _err
	}
	err := xml.Unmarshal(content, &cfg)
	return &cfg, err
}

// 转换列与字段的映射
func parseColumnMapping(str string) []ColumnMapping {
	re, err := regexp.Compile("([^:]+):([^;]*);*\\s*")
	if err != nil {
		return nil
	}
	var matches [][]string = re.FindAllStringSubmatch(str, -1)
	if matches == nil {
		return nil
	}
	columnsMapping := make([]ColumnMapping, len(matches))
	for i, v := range matches {
		columnsMapping[i] = ColumnMapping{Field: v[1], Name: v[2]}
	}
	return columnsMapping
}

func Export(portal IDataExportPortal, parameters *Params,
	provider IDataExportProvider) []byte {
	rows, _, _ := portal.GetSchemaAndData(parameters.Params)
	names := portal.GetExportColumnNames(
		parameters.ExportColumnNames)
	return provider.Export(rows, parameters.ExportColumnNames, names)
}

func GetExportParams(paramMappings string, columnNames []string) *Params {
	parameters := make(map[string]string)
	if paramMappings != "" {
		paramMappings = strings.Replace(paramMappings,
			"%3d", "=", -1)
		var paramsArr, splitArr []string
		paramsArr = strings.Split(paramMappings, ";")
		//添加传入的参数
		for _, v := range paramsArr {
			splitArr = strings.Split(v, ":")
			parameters[splitArr[0]] = v[len(splitArr[0])+1:]
		}
	}
	return &Params{ExportColumnNames: columnNames, Params: parameters}

}

// 格式化sql语句
func SqlFormat(sql string, ht map[string]string) (formatted string) {
	formatted = sql
	for k, v := range ht {
		formatted = strings.Replace(formatted, "{"+k+"}", v, -1)
	}
	return formatted
}
