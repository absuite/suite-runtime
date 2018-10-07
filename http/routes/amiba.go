package routes

import (
	"github.com/absuite/suite-runtime/http/controllers/amiba"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/absuite/suite-runtime/services/amiba"

	"github.com/kataras/iris/mvc"
)

func registerAmiba(c *Context) {
	//model route
	repo := repositories.NewModelRepo(c.Orm)
	modelSv := amibaServices.NewModelSv(repo)
	modeling := mvc.New(c.App.Party("/api/amiba/models"))
	modeling.Register(modelSv)
	modeling.Handle(new(amibaControllers.ModelController))

	//price route
	priceSv := amibaServices.NewPricelSv(repo)
	price := mvc.New(c.App.Party("/api/amiba/prices"))
	price.Register(priceSv)
	price.Handle(new(amibaControllers.PriceController))
}
