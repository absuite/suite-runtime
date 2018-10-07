package middleware

import (
	"github.com/kataras/iris"
)

func Ent(ctx iris.Context) {
	ctx.Values().Set("Ent", ctx.GetHeader("Ent"))
	ctx.Next()
}
