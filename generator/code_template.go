package generator

import (
	"regexp"
	"strings"
)

var(
	predefineRegexp = regexp.MustCompilePOSIX("\\!([^:-]+):(.+)")
)
type CodeTemplate struct{
	template string
	predefine map[string]string
}

func NewTemplate(s string)*CodeTemplate{
	return (&CodeTemplate{}).configure(s)
}

func (g *CodeTemplate) configure(s string)*CodeTemplate {
	g.predefine = make(map[string]string)
	for _,match := range predefineRegexp.FindAllStringSubmatch(s,-1){
		g.predefine[match[1]]= match[2]
	}
	g.template = predefineRegexp.ReplaceAllString(s,"")
	return g
}

//  获取模板内容
func (g *CodeTemplate) String() string {
	return g.template
}

// 获取预定义的参数
func (g *CodeTemplate) Predefine(key string)(string,bool){
	n,ok := g.predefine[key]
	return n,ok
}

func (g *CodeTemplate) Replace(s, old string, n int)*CodeTemplate {
	g.template = strings.Replace(g.template, s, old, n)
	return g
}

func resolveRepTag(g *CodeTemplate)*CodeTemplate {
	return g.Replace("<Ptr>", "{{.Ptr}}", -1).
		Replace("<E>", "{{.E}}", -1).
		Replace("<E2>", "{{.E2}}", -1).
		Replace("<R>", "{{.R}}", -1).
		Replace("<R2>", "{{.R2}}", -1).
		Replace("<IsPK>", "{{.IsPK}}", -1)
}

func init() {
	TPL_ENTITY_REP = resolveRepTag(TPL_ENTITY_REP)
	TPL_ENTITY_REP_INTERFACE = resolveRepTag(TPL_ENTITY_REP_INTERFACE)
	TPL_REPO_FACTORY = resolveRepTag(TPL_REPO_FACTORY)
}


