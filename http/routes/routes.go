package routes

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
)

type Context struct {
	Orm *xorm.Engine
	App *iris.Application
}

func Register(r *Context) {
	registerAmiba(r)
}
