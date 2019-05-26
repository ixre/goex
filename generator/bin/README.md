# 代码生成器

使用Go编写的代码生成器,可根据模板定制生成代码.

特点:
- 支持多mysql和postgresql数据库
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