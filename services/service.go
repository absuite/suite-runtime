package services

import (
	"github.com/absuite/suite-runtime/services/amiba"
	"github.com/absuite/suite-runtime/services/cbo"
	"github.com/ggoop/goutils/glog"
	"github.com/kataras/iris/hero"
	"go.uber.org/dig"
)

func afterRegister(container *dig.Container) {
	if err := container.Invoke(func(s1 cboServices.EntSv, s2 cboServices.OrgSv, s3 cboServices.PeriodSv, s4 cboServices.DeptSv) {
		hero.Register(s1)
		hero.Register(s2)
		hero.Register(s3)
		hero.Register(s4)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if err := container.Invoke(func(s1 cboServices.ItemSv, s2 cboServices.WhSv, s3 cboServices.ProjectSv, s4 cboServices.DocTypeSv, s5 cboServices.TraderSv) {
		hero.Register(s1)
		hero.Register(s2)
		hero.Register(s3)
		hero.Register(s4)
		hero.Register(s5)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}

	if err := container.Invoke(func(s1 amibaServices.PurposeSv, s2 amibaServices.GroupSv, s3 amibaServices.PricelSv, s4 amibaServices.ModelSv) {
		hero.Register(s1)
		hero.Register(s2)
		hero.Register(s3)
		hero.Register(s4)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
func Register(container *dig.Container) {
	//企业
	if err := container.Provide(cboServices.NewEntSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//组织服务
	if err := container.Provide(cboServices.NewOrgSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//期间服务
	if err := container.Provide(cboServices.NewPeriodSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//部门服务
	if err := container.Provide(cboServices.NewDeptSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//物料服务
	if err := container.Provide(cboServices.NewItemSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//存储地点服务
	if err := container.Provide(cboServices.NewWhSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//项目服务
	if err := container.Provide(cboServices.NewProjectSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//客商服务
	if err := container.Provide(cboServices.NewTraderSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//单据类型服务
	if err := container.Provide(cboServices.NewDocTypeSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//amiba
	// 核算目的
	if err := container.Provide(amibaServices.NewPurposeSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	// 价表
	if err := container.Provide(amibaServices.NewGroupSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//阿米巴模型服务
	if err := container.Provide(amibaServices.NewModelSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//价表服务
	if err := container.Provide(amibaServices.NewPricelSv); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}

	afterRegister(container)
}
