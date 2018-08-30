package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/jsix/goex/generator"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/shell"
	"github.com/jsix/gof/web/form"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	var version = "1.0.5"
	var genDir string   //输出目录
	var confPath string //设置目录
	var tplDir string   //模板目录
	var table string
	var arch string //代码架构
	var debug bool
	var printVer bool

	flag.StringVar(&genDir, "out", "./output", "path of output directory")
	flag.StringVar(&tplDir, "tpl", "./templates", "path of code templates directory")
	flag.StringVar(&confPath, "conf", "./", "config path")
	flag.StringVar(&table, "table", "", "table name or table prefix")
	flag.StringVar(&arch, "arch", "", "program language")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&printVer, "v", false, "print version")
	flag.Parse()
	if printVer {
		fmt.Println("GofGenerator v" + version)
		return
	}
	registry, err := gof.NewRegistry(confPath, ".")
	if err != nil {
		log.Println("[ Gen][ Fail]:", err.Error())
		return
	}
	defer crashRecover()
	re := registry.Use("gen")
	// 获取bash启动脚本，默认unix系统包含了bash，windows下需指定
	bashExec := ""
	if runtime.GOOS == "windows" {
		bashExec = re.GetString("command.win_bash")
	}
	// 生成之前执行操作
	if err := runBefore(re, bashExec); err != nil {
		log.Fatalln("[ Gen][ Fail]:", err.Error())
	}
	// 初始化生成器
	dbName := re.GetString("database.name")
	d := &orm.MySqlDialect{}
	ds := orm.DialectSession(getDb(re), d)
	dg := generator.DBCodeGenerator()
	dg.IdUpper = true
	// 获取表格并转换
	tables, err := dg.ParseTables(ds.TablesByPrefix(dbName, table))
	if err != nil {
		log.Println("[ Gen][ Fail]:", err.Error())
		return
	}
	// 生成代码
	if err := genByArch(arch, dg, tables, genDir, tplDir); err != nil {
		log.Fatalln("[ Gen][ Fail]:", err.Error())
	}
	// 生成之后执行操作
	if err := runAfter(re, bashExec); err != nil {
		log.Fatalln("[ Gen][ Fail]:", err.Error())
	}
	log.Println("[ Gen][ Success]: generate successfully!")
}

func runBefore(re *gof.RegistryTree, bashExec string) error {
	beforeRun := strings.TrimSpace(re.GetString("command.before"))
	return execCommand(beforeRun, bashExec)
}

func runAfter(re *gof.RegistryTree, bashExec string) error {
	afterRun := strings.TrimSpace(re.GetString("command.after"))
	return execCommand(afterRun, bashExec)
}

// 执行命令
func execCommand(command string, bashExec string) error {
	// 生成之后执行操作
	if command != "" {
		if bashExec != "" {
			if strings.Contains(bashExec, " ") {
				bashExec = "\"" + bashExec + "\""
			}
			command = bashExec + " " + command
		}
		_, _, err := shell.StdRun(command)
		return err
	}
	return nil
}

// 根据规则生成代码
func genByArch(arch string, dg *generator.Session, tables []*generator.Table,
	genDir string, tplDir string) error {
	// 生成代码
	switch arch {
	case "repo":
		return genGoCode(dg, tables, genDir, tplDir)
	default:
		return genCode(dg, tables, genDir, tplDir)
	}
	return nil
}

// 获取数据库连接
func getDb(r *gof.RegistryTree) *sql.DB {
	//数据库连接字符串
	//root@tcp(127.0.0.1:3306)/db_name?charset=utf8
	var prefix = "database"
	driver := r.GetString(prefix + ".driver")
	dbCharset := r.GetString(prefix + ".charset")
	if dbCharset == "" {
		dbCharset = "utf8"
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=Local",
		r.GetString(prefix+".user"),
		r.GetString(prefix+".pwd"),
		r.GetString(prefix+".server"),
		r.Get(prefix + ".port").(int64),
		r.GetString(prefix+".name"),
		dbCharset,
	)
	db, err := sql.Open(driver, connStr)
	if err == nil {
		db.SetMaxIdleConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Second * 10)
		err = db.Ping()
	}
	if err != nil {
		defer db.Close()
		panic(err)
		return nil
	}
	return db
}

// 恢复应用
func crashRecover() {
	r := recover()
	if r != nil {
		fmt.Println(fmt.Sprintf("[ Gen][ Error]: %v", r))
	}
}

// 生成代码
func genCode(s *generator.Session, tables []*generator.Table, genDir string, tplDir string) error {
	tplMap := map[string]generator.CodeTemplate{}
	sliceSize := len(tplDir) - 1
	if tplDir[sliceSize] == '/' {
		tplDir = tplDir + "/"
		sliceSize += 1
	}
	err := filepath.Walk(tplDir, func(path string, info os.FileInfo, err error) error {
		// 如果模板名称以"_"开头，则忽略
		if !info.IsDir() && info.Name()[0] != '_' {
			tp, err := s.ParseTemplate(path)
			if err == nil {
				tplMap[path[sliceSize:]] = tp
			}
			return errors.New("template:" + info.Name() + "-" + err.Error())
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(tplMap) == 0 {
		return errors.New("no any code template")
	}
	for _, tb := range tables {
		for path, tpl := range tplMap {
			str := s.GenerateCode(tb, tpl, "", true, "")
			generator.SaveFile(str, genDir+"/"+joinFilePath(path, tb.Name))
		}
	}
	return err
}

// 连接文件路径
func joinFilePath(path string, tableName string) string {
	i := strings.Index(path, ".")
	if i != -1 {
		return strings.Join([]string{path[:i], "_", tableName, ".", path[i+1:]}, "")
	}
	return path + tableName
}

// 生成Go代码
func genGoCode(dg *generator.Session, tables []*generator.Table,
	genDir string, tplDir string) error {
	// 设置变量
	dg.Var(generator.VModelPkg, "xupms/src/model")
	dg.Var(generator.VModelPkgName, "model")
	dg.Var(generator.VIRepoPkg, "xupms/src/repo")
	dg.Var(generator.VIRepoPkgName, "repo")
	dg.Var(generator.VRepoPkg, "xupms/src/repo")
	dg.Var(generator.VRepoPkgName, "repo")
	// 读取自定义模板
	listTP, _ := dg.ParseTemplate(tplDir + "/grid_list.html")
	editTP, _ := dg.ParseTemplate(tplDir + "/entity_edit.html")
	ctrTpl, _ := dg.ParseTemplate(tplDir + "/entity.html")
	var err error
	// 初始化表单引擎
	fe := &form.Engine{}
	for _, tb := range tables {
		entityPath := genDir + "model/" + tb.Name + ".go"
		iRepPath := genDir + "repo/auto_iface_" + tb.Name + "_repo.go"
		repPath := genDir + "repo/auto_" + tb.Name + "_repo.go"
		dslPath := genDir + "form/" + tb.Name + ".form"
		htmPath := genDir + "html/" + tb.Name + ".html"
		//生成实体
		str := dg.TableToGoStruct(tb)
		generator.SaveFile(str, entityPath)
		//生成仓储结构
		str = dg.TableToGoRepo(tb, true, "model.")
		generator.SaveFile(str, repPath)
		//生成仓储接口
		str = dg.TableToGoIRepo(tb, true, "")
		generator.SaveFile(str, iRepPath)
		//生成表单DSL
		f := fe.TableToForm(tb.Raw)
		err = fe.SaveDSL(f, dslPath)
		//生成表单
		if err == nil {
			_, err = fe.SaveHtmlForm(f, form.TDefaultFormHtml, htmPath)
		}
		if err != nil {
			return err
		}
		// 生成列表文件
		str = dg.GenerateCode(tb, listTP, "", true, "")
		generator.SaveFile(str, genDir+"html_list/"+tb.Name+"_list.html")
		// 生成表单文件
		str = dg.GenerateCode(tb, editTP, "", true, "")
		generator.SaveFile(str, genDir+"html_edit/"+tb.Name+"_edit.html")
		// 生成控制器
		str = dg.GenerateCode(tb, ctrTpl, "", true, "")
		generator.SaveFile(str, genDir+"c/"+tb.Name+"_c.go")
	}
	// 生成仓储工厂
	code := dg.GenerateTablesCode(tables, generator.TPL_REPO_FACTORY)
	generator.SaveFile(code, genDir+"repo/auto_repo_factory.go")
	//格式化代码
	shell.Run("gofmt -w " + genDir)
	return err
}
