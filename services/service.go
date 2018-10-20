package services

import (
	"github.com/absuite/suite-runtime/services/amiba"
	"github.com/absuite/suite-runtime/services/cbo"
	"github.com/ggoop/goutils/glog"
	"go.uber.org/dig"
)

func afterRegister(container *dig.Container) {
	if err := container.Invoke(func(l amibaServices.ModelSv) {
		glog.Errorf("di Invoke success:%s")
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if err := container.Invoke(func(l amibaServices.PricelSv) {
		glog.Errorf("di Invoke success:%s")
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
func Register(container *dig.Container) {

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
	//单据类型服务
	if err := container.Provide(cboServices.NewDocTypeSv); err != nil {
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
}
