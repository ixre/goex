package echox

import (
	"github.com/labstack/echo/v4"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : group
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-12-03 11:42
 * description :
 * history :
 */

type GroupHandler interface {
	Handle(group *echo.Group)
}
