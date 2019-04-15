package generator

import (
	"github.com/ixre/gof/util"
)

// 保存到文件
func SaveFile(s string, path string) error {
	return util.BytesToFile([]byte(s), path)
}
