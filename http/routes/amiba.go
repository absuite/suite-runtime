package routes

import (
	"github.com/absuite/suite-runtime/http/controllers/amiba"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/absuite/suite-runtime/services/amiba"
	"github.com/ggoop/goutils/glog"

	"github.com/kataras/iris/mvc"
)

func registerAmiba(c *Context) {
	//model route
	repo := repositories.NewModelRepo(c.Orm)
	modelSv := amibaServices.NewModelSv(repo)
	modeling := mvc.New(c.App.Party("/api/amiba/models"))
	modeling.Register(modelSv)
	modeling.Handle(new(amibaControllers.ModelController))
	glog.Printf("注册 %s 路由成功", "/api/amiba/models")

	//price route
	priceSv := amibaServices.NewPricelSv(repo)
	priceSv.CacheAll() //缓存价表数据
	price := mvc.New(c.App.Party("/api/amiba/prices"))
	price.Register(priceSv)
	price.Handle(new(amibaControllers.PriceController))
	glog.Printf("注册 %s 路由成功", "/api/amiba/prices")
}
