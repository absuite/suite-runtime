package routes

import (
	"github.com/absuite/suite-runtime/http/controllers/amiba"
	"github.com/ggoop/goutils/glog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"go.uber.org/dig"
)

func registerAmiba(container *dig.Container) {
	//model route
	if err := container.Invoke(func(app *iris.Application) {
		m := mvc.New(app.Party("/api/amiba/models"))
		m.Handle(new(amibaControllers.ModelController))
		glog.Printf("注册 %s 路由成功", "/api/amiba/models")
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//price route
	if err := container.Invoke(func(app *iris.Application) {
		price := mvc.New(app.Party("/api/amiba/prices"))
		price.Handle(new(amibaControllers.PriceController))
		glog.Printf("注册 %s 路由成功", "/api/amiba/prices")
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
