package routes

import (
	"github.com/absuite/suite-runtime/http/controllers/amiba"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/absuite/suite-runtime/services/amiba"
	"github.com/absuite/suite-runtime/services/cbo"
	"github.com/ggoop/goutils/glog"
	"github.com/kataras/iris/hero"
	"github.com/kataras/iris/mvc"
)

func registerAmiba(c *Context) {
	repo := repositories.NewModelRepo(c.Orm)
	//组织服务
	hero.Register(cboServices.NewOrgSv(repo))
	//期间服务
	hero.Register(cboServices.NewPeriodSv(repo))
	//部门服务
	hero.Register(cboServices.NewDeptSv(repo))
	//物料服务
	hero.Register(cboServices.NewItemSv(repo))
	//存储地点服务
	hero.Register(cboServices.NewWhSv(repo))
	//项目服务
	hero.Register(cboServices.NewWhSv(repo))
	//单据类型服务
	hero.Register(cboServices.NewDocTypeSv(repo))

	//阿米巴模型服务
	hero.Register(amibaServices.NewModelSv(repo))
	//价表服务
	hero.Register(amibaServices.NewPricelSv(repo))

	//model route

	modeling := mvc.New(c.App.Party("/api/amiba/models"))
	modeling.Handle(new(amibaControllers.ModelController))
	glog.Printf("注册 %s 路由成功", "/api/amiba/models")

	//price route

	price := mvc.New(c.App.Party("/api/amiba/prices"))
	price.Handle(new(amibaControllers.PriceController))
	glog.Printf("注册 %s 路由成功", "/api/amiba/prices")
}
