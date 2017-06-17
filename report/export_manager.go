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
	fileSize      int64
	watch         bool
	columnMapping []ColumnMapping
	sqlConfig     *ItemConfig
	Base          *DataExportPortal
	PortalKey     string
	//管理器
	ItemManager *ExportItemManager
}

func (e *ExportItem) formatMappingString(str string) string {
	reg := regexp.MustCompile("\\s|\\n")
	return reg.ReplaceAllString(e.sqlConfig.ColumnMapping, "")
}

//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
func (e *ExportItem) GetColumnNames() []ColumnMapping {
	if e.columnMapping == nil {
		e.sqlConfig.ColumnMapping = e.formatMappingString(e.sqlConfig.ColumnMapping)
		e.columnMapping = parseColumnMapping(e.sqlConfig.ColumnMapping)
	}
	return e.columnMapping
}

func (e *ExportItem) GetExportColumnIndexAndName(
	exportColumnNames []string) (dict map[string]string) {
	dict = make(map[string]string)
	for _, cName := range exportColumnNames {
		for _, cMap := range e.GetColumnNames() {
			if cMap.Name == cName {
				dict[cMap.Field] = cMap.Name
				break
			}
		}
	}
	return dict
}

// 检查SQL配置
func (portal *ExportItem) checkSqlConfig() (err error) {
	if portal.sqlConfig == nil {
		dir, _ := os.Getwd()
		portal.sqlConfig, err = readItemConfigFromXml(
			strings.Join([]string{dir, portal.ItemManager.RootPath,
				portal.PortalKey, ".xml"}, ""))
		if err != nil {
			portal.sqlConfig = nil
			return err
		}
	}
	return nil
}

//获取统计数据
func (portal *ExportItem) GetTotalView(ht map[string]string) (row map[string]interface{}) {
	return nil
}

func (portal *ExportItem) GetSchemaAndData(ht map[string]string) (rows []map[string]interface{}, total int, err error) {
	if err := portal.checkSqlConfig(); err != nil {
		return nil, 0, err
	}

	total = 0
	var _rows *sql.Rows
	_db := portal.ItemManager.DbGetter.GetDB()

	//初始化添加参数
	if _, e := ht["pageSize"]; !e {
		ht["pageSize"] = "10000000000"
	}
	if _, e := ht["pageIndex"]; !e {
		ht["pageIndex"] = "1"
	}

	pi, _ := ht["pageIndex"]
	ps, _ := ht["pageSize"]
	pageIndex, _ := strconv.Atoi(pi)
	pageSize, _ := strconv.Atoi(ps)

	if pageIndex > 0 {
		ht["page_start"] = strconv.Itoa((pageIndex - 1) * pageSize)
	} else {
		ht["page_start"] = "0"
	}
	ht["page_end"] = strconv.Itoa(pageIndex * pageSize)
	ht["page_size"] = strconv.Itoa(pageSize)

	//统计总行数
	if portal.sqlConfig.Total != "" {
		sql := SqlFormat(portal.sqlConfig.Total, ht)
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
	if portal.sqlConfig.Query != "" {
		sql := SqlFormat(portal.sqlConfig.Query, ht)
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
func (portal *ExportItem) GetJsonData(ht map[string]string) string {
	result, err := json.Marshal(nil)
	if err != nil {
		return "{error:'" + err.Error() + "'}"
	}
	return string(result)
}

//导出项管理器
type ExportItemManager struct {
	watch bool
	//配置存放路径
	RootPath string
	//配置扩展名
	CfgFileExt string
	//数据库连接
	DbGetter IDbProvider //接口类型不需要加*
	//导出项集合
	exportItems map[string]*ExportItem
}

func NewExportManager(db IDbProvider, watch bool) *ExportItemManager {
	return &ExportItemManager{
		watch:       watch,
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
			manager.DbGetter, manager.watch)
		//manager.exportItems[portalKey] = item
	}
	return item
}

// 创建导出项,watch：是否监视文件变化
func (manager *ExportItemManager) loadExportItem(portalKey string,
	dbp IDbProvider, watch bool) *ExportItem {
	dir, _ := os.Getwd()
	arr := []string{dir, manager.RootPath, portalKey, manager.CfgFileExt}
	filePath := strings.Join(arr, "")
	if f, err := os.Stat(filePath); err == nil && f.IsDir() == false {
		cfg, err := readItemConfigFromXml(filePath)
		if err == nil {
			return &ExportItem{
				fileSize:    f.Size(),
				watch:       manager.watch,
				sqlConfig:   cfg,
				PortalKey:   portalKey,
				ItemManager: manager,
			}
		}
	}
	return nil
}
