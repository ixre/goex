package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/jsix/goex/generator"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/shell"
	"github.com/jsix/gof/web/form"
	"time"
)

func main() {
	var genDir string   //输出目录
	var confPath string //设置目录
	flag.StringVar(&genDir, "out", "generated_code/", "generated code output directory")
	flag.StringVar(&confPath, "conf", "./", "config path")
	flag.Parse()

	registry, err := gof.NewRegistry(confPath, ".")
	if err != nil {
		panic(err)
	}
	// 初始化生成器
	d := &orm.MySqlDialect{}
	ds := orm.DialectSession(getDb(registry), d)
	dg := generator.DBCodeGenerator()
	dg.IdUpper = true
	// 获取表格并转换
	tables, err := dg.ParseTables(ds.TablesByPrefix("", "uams_"))
	if err != nil {
		return
	}
	// 设置变量
	dg.Var(generator.V_ModelPkg, "xupms/src/model")
	dg.Var(generator.V_ModelPkgName, "model")
	dg.Var(generator.V_IRepoPkg, "xupms/src/repo")
	dg.Var(generator.V_IRepoPkgName, "repo")
	dg.Var(generator.V_RepoPkg, "xupms/src/repo")
	dg.Var(generator.V_RepoPkgName, "repo")
	// 读取自定义模板
	listTP, _ := dg.ParseTemplate("code_templates/grid_list.html")
	editTP, _ := dg.ParseTemplate("code_templates/entity_edit_table.html")
	ctrTpl, _ := dg.ParseTemplate("code_templates/entity_c._go")
	// 初始化表单引擎
	fe := &form.Engine{}
	for _, tb := range tables {
		entityPath := genDir + "model/" + tb.Name + ".go"
		iRepPath := genDir + "repo/auto_i" + tb.Name + "_repo.go"
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

	//t.Log("生成成功")
	log.Println("----")
}

// 获取数据库连接
func getDb(r *gof.Registry) *sql.DB {
	//数据库连接字符串
	//root@tcp(127.0.0.1:3306)/db_name?charset=utf8
	var prefix = "gen.database"
	driver := r.GetString(prefix + ".driver")
	dbCharset := r.GetString(prefix + ".charset")
	if dbCharset == "" {
		dbCharset = "utf8"
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local",
		r.GetString(prefix+".user"),
		r.GetString(prefix+".pwd"),
		r.GetString(prefix+".server"),
		r.Get(prefix+".port").(int64),
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
		log.Println("[ Gof][ MySQL] " + err.Error())
		return nil
	}
	return db
}
