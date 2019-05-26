# 代码生成器

使用Go编写的代码生成器,可根据模板定制生成代码.

特点:
- 支持mysql和postgresql数据库
- 支持Go,JAVA,Kotlin,Html,C#语言
- 支持代码模板, 支持模板函数

资源:
- [下载地址](https://github.com/ixre/goex/releases/)
- [Go模板语法-中](http://www.g-var.com/posts/translation/hugo/hugo-21-go-template-primer/)
- [Go模板语法-English](https://golang.org/pkg/text/template/)

## 快速开始

1. 配置数据源
```
下载安装包,解压修改gen.conf文件进行数据源配置.
```
2. 定制修改模板
```
根据实际需求对模板进行修改, 或创建自己的模板. 模板语法请参考: Go Template
```
3. 运行命令生成代码
```bash
gof-gen -conf gen.conf
```

## 预定义语法

预定义语法用来在代码模板中定义一些数据, 在生成代码时预定义语法不输入任何内容.
预定义语法格式为: !预定义参数名:预定义参数值

目前,支持的预定义语法如下:

- !target : 用来定义代码文件存放的目标路径

## 模板

模板主要包含以下对象: 

- global
- table
- columns


### global

> 用于读取全局变量, global的属性均以大写开头; global为小写.

输出生成器的版本号
```
// this file created by generate {{.global.Version}}
```
输出包名,包名通过配置文件配置.格式为: com/pkg
```
package {{.global.Pkg}}
```
如果是Java或其他语言, 包名以"."分割, 可使用pkg函数,如:
```
// java package
package {{pkg "java" .global.Pkg}}
// c# namespace
namespace {{pkg "csharp" .global.Pkg}}
```

## 模板示例

以下代码用于生成Java的Pojo对象, 更多示例点击[这里](./templates)

```
!target:{{.global.Pkg}}/pojo/{{.table.Title}}Entity.java
package {{pkg "java" .global.Pkg}}.pojo;

import javax.persistence.Basic;
import javax.persistence.Id;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.persistence.GenerationType;
import javax.persistence.GeneratedValue;

/** {{.table.Comment}} */
@Entity
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
public class {{.table.Title}}Entity {
    {{range $i,$c := .columns}}{{$type := type "java" $c.TypeId}}
    private {{$type}} {{$c.Name}}
    public void set{{$c.Title}}({{$type}} {{$c.Name}}){
        this.{{$c.Name}} = {{$c.Name}}
    }

    /** {{$c.Comment}} */{{if $c.IsPk}}
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"{{if not $c.NotNull}}, nullable = true{{end}} {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    public {{$type}} get{{$c.Title}}() {
        return this.{{$c.Name}};
    }
    {{end}}
}

```