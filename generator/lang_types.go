package generator

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)

func JavaTypes(typeId int) string {
	switch typeId {
	case orm.GoTypeBoolean:
		return "Boolean"
	case orm.GoTypeInt64:
		return "Long"
	case orm.GoTypeFloat32:
		return "Float"
	case orm.GoTypeFloat64:
		return "Double"
	case orm.GoTypeInt32:
		return "int"
	case orm.GoTypeString:
		return "String"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
