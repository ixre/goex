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
	"encoding/json"
	"errors"
	"github.com/jsix/gof/db"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var _ IDataExportPortal = new(ExportItem)

// 导出项目
type ExportItem struct {
	columnMapping []ColumnMapping
	sqlConfig     *ItemConfig
	dbProvider    IDbProvider
	PortalKey     string
}

func (e *ExportItem) formatMappingString(str string) string {
	reg := regexp.MustCompile("\\s|\\n")
	return reg.ReplaceAllString(e.sqlConfig.ColumnMapping, "")
}

//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
func (e *ExportItem) GetColumnMapping() []ColumnMapping {
	if e.columnMapping == nil {
		e.sqlConfig.ColumnMapping = e.formatMappingString(e.sqlConfig.ColumnMapping)
		e.columnMapping = parseColumnMapping(e.sqlConfig.ColumnMapping)
	}
	return e.columnMapping
}

func (e *ExportItem) GetExportColumnNames(
	exportColumns []string) (names []string) {
	names = make([]string, len(exportColumns))
	mapping := e.GetColumnMapping()
	for i, cName := range exportColumns {
		for _, cMap := range mapping {
			if cMap.Field == cName {
				names[i] = cMap.Name
				break
			}
		}
	}
	return names
}

//获取统计数据
func (e *ExportItem) GetTotalView(ht map[string]string) (row map[string]interface{}) {
	return nil
}

func (e *ExportItem) GetSchemaAndData(p Params) (rows []map[string]interface{}, total int, err error) {
	if e == nil || e.dbProvider == nil {
		return nil, 0, errors.New("no such export item")
	}
	total = 0
	var _rows *sql.Rows
	_db := e.dbProvider.GetDB()

	//初始化添加参数
	if _, e := p["pageSize"]; !e {
		p["pageSize"] = "10000000000"
	}
	if _, e := p["pageIndex"]; !e {
		p["pageIndex"] = "1"
	}

	pi, _ := p["pageIndex"]
	ps, _ := p["pageSize"]
	pageIndex, _ := strconv.Atoi(pi)
	pageSize, _ := strconv.Atoi(ps)

	if pageIndex > 0 {
		p["page_start"] = strconv.Itoa((pageIndex - 1) * pageSize)
	} else {
		p["page_start"] = "0"
	}
	p["page_end"] = strconv.Itoa(pageIndex * pageSize)
	p["page_size"] = strconv.Itoa(pageSize)

	//统计总行数
	if e.sqlConfig.Total != "" {
		sql := SqlFormat(e.sqlConfig.Total, p)
		smt, err := _db.Prepare(sql)

		if err != nil {
			log.Println("[ Export][ Error] -", err.Error(), "\n", sql)
			return nil, 0, err
		}

		row := smt.QueryRow()
		if row != nil {
			err = row.Scan(&total)
			if err != nil {
				log.Println("[ Export][ Error] -", err.Error(), "\n", sql)
				return nil, total, err
			}
		}
	}

	//获得数据
	if e.sqlConfig.Query != "" {
		sql := SqlFormat(e.sqlConfig.Query, p)
		//log.Println("-----",sql)
		sqlLines := strings.Split(sql, ";\n")
		if t := len(sqlLines); t > 1 {
			for i, v := range sqlLines {
				if i != t-1 {
					if smt, err := _db.Prepare(v); err == nil {
						smt.Exec()
					}
				}
			}
			sql = sqlLines[t-1]

		}

		smt, err := _db.Prepare(sql)

		if err != nil {
			log.Println("[ Export][ Error] -", err.Error(), "\n", sql)
			return nil, total, err
		}
		_rows, err = smt.Query()
		if err != nil {
			log.Println("[ Export][ Error] -", err.Error(), "\n", sql)
			return nil, total, err
		}
		defer _rows.Close()
	}

	return db.RowsToMarshalMap(_rows), total, err
}

//获取要导出的数据Json格式
func (e *ExportItem) GetJsonData(ht map[string]string) string {
	result, err := json.Marshal(nil)
	if err != nil {
		return "{error:'" + err.Error() + "'}"
	}
	return string(result)
}

func (e *ExportItem) Export(parameters *ExportParams,
	provider IExportProvider, formatter IExportFormatter) []byte {
	rows, _, _ := e.GetSchemaAndData(parameters.Params)
	names := e.GetExportColumnNames(
		parameters.ExportFields)
	fmtArray := []IExportFormatter{interFmt}
	if formatter != nil {
		fmtArray = append(fmtArray, formatter)
	}
	return provider.Export(rows, parameters.ExportFields, names, fmtArray)
}

//导出项管理器
type ExportItemManager struct {
	//配置存放路径
	RootPath string
	//配置扩展名
	CfgFileExt string
	//数据库连接
	DbGetter IDbProvider //接口类型不需要加*
	//导出项集合
	exportItems map[string]*ExportItem
}

func NewExportManager(db IDbProvider) *ExportItemManager {
	return &ExportItemManager{
		DbGetter:    db,
		RootPath:    "/conf/query/",
		CfgFileExt:  ".xml",
		exportItems: make(map[string]*ExportItem),
	}
}

//获取导出项
func (manager *ExportItemManager) GetExportItem(portalKey string) IDataExportPortal {
	item, exist := manager.exportItems[portalKey]
	if !exist {
		item = manager.loadExportItem(portalKey,
			manager.DbGetter)
		if !WATCH_CONF_FILE {
			manager.exportItems[portalKey] = item
		}
	}
	return item
}

// 创建导出项,watch：是否监视文件变化
func (manager *ExportItemManager) loadExportItem(portalKey string,
	dbp IDbProvider) *ExportItem {
	dir, _ := os.Getwd()
	arr := []string{dir, manager.RootPath, portalKey, manager.CfgFileExt}
	filePath := strings.Join(arr, "")
	f, err := os.Stat(filePath)
	if err == nil && f.IsDir() == false {
		cfg, err1 := readItemConfigFromXml(filePath)
		if err1 == nil {
			return &ExportItem{
				sqlConfig:  cfg,
				PortalKey:  portalKey,
				dbProvider: dbp,
			}
		}
		err = err1
	}
	if err != nil {
		log.Println("[ Export][ Error]:", err.Error(), "; portal:", portalKey)
	}
	return nil
}
