package amibaServices

import (
	"time"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/models/cbo"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func (s *modelSv) DtiLog_Init(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelId string) {
	dtiModeling := amibaModels.DtiModeling{EntId: ent.Id, PurposeId: purpose.ID, PeriodId: period.Id, ModelId: modelId}
	if _, err := s.repo.Get(&dtiModeling); err != nil {
		glog.Errorf("企业:%v,核算:%v,期间:%v,获取模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
	}
	dtiModeling.Succeed = 0
	dtiModeling.Msg = "准备开始..."

	if dtiModeling.Id == "" {
		dtiModeling.Id = utils.GUID()
		if _, err := s.repo.Insert(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,创建模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	} else {
		if _, err := s.repo.Id(dtiModeling.Id).Cols("start_time", "end_time", "Msg", "Succeed").Update(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,更新模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	}
}

func (s *modelSv) DtiLog_Begin(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelId string) {
	dtiModeling := amibaModels.DtiModeling{EntId: ent.Id, PurposeId: purpose.ID, PeriodId: period.Id, ModelId: modelId}
	if _, err := s.repo.Get(&dtiModeling); err != nil {
		glog.Errorf("企业:%v,核算:%v,期间:%v,获取模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
	}
	dtiModeling.StartTime = time.Now()
	dtiModeling.Succeed = 0
	dtiModeling.Msg = "开始"

	if dtiModeling.Id == "" {
		dtiModeling.Id = utils.GUID()
		if _, err := s.repo.Insert(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,创建模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	} else {
		if _, err := s.repo.Id(dtiModeling.Id).Cols("start_time", "end_time", "Msg", "Succeed").Update(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,更新模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	}
}
func (s *modelSv) DtiLog_Log(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelId string, msg string) {
	dtiModeling := amibaModels.DtiModeling{EntId: ent.Id, PurposeId: purpose.ID, PeriodId: period.Id, ModelId: modelId}
	if _, err := s.repo.Get(&dtiModeling); err != nil {
		glog.Errorf("企业:%v,核算:%v,期间:%v,获取模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
	}
	dtiModeling.Msg = msg
	if dtiModeling.Id == "" {
		dtiModeling.Id = utils.GUID()
		if _, err := s.repo.Insert(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,创建模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	} else {
		if _, err := s.repo.Id(dtiModeling.Id).Cols("Msg").Update(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,更新模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	}
}
func (s *modelSv) DtiLog_Error(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelId string, err error) {
	dtiModeling := amibaModels.DtiModeling{EntId: ent.Id, PurposeId: purpose.ID, PeriodId: period.Id, ModelId: modelId}
	if _, err := s.repo.Get(&dtiModeling); err != nil {
		glog.Errorf("企业:%v,核算:%v,期间:%v,获取模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
	}
	dtiModeling.Msg = err.Error()
	dtiModeling.EndTime = time.Now()
	dtiModeling.Succeed = 0
	if dtiModeling.Id == "" {
		dtiModeling.Id = utils.GUID()
		if _, err := s.repo.Insert(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,创建模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	} else {
		if _, err := s.repo.Id(dtiModeling.Id).Cols("Msg", "end_time", "Succeed").Update(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,更新模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	}
}
func (s *modelSv) DtiLog_Success(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelId string) {
	dtiModeling := amibaModels.DtiModeling{EntId: ent.Id, PurposeId: purpose.ID, PeriodId: period.Id, ModelId: modelId}
	if _, err := s.repo.Get(&dtiModeling); err != nil {
		glog.Errorf("企业:%v,核算:%v,期间:%v,获取模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
	}
	dtiModeling.EndTime = time.Now()
	dtiModeling.Succeed = 1
	dtiModeling.Msg = "完成"
	if dtiModeling.Id == "" {
		dtiModeling.Id = utils.GUID()
		if _, err := s.repo.Insert(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,创建模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	} else {
		if _, err := s.repo.Id(dtiModeling.Id).Cols("end_time", "Succeed", "Msg").Update(&dtiModeling); err != nil {
			glog.Errorf("企业:%v,核算:%v,期间:%v,更新模型日志错误,%s", ent.Name, purpose.Name, period.Name, err)
		}
	}
}

func (s *modelSv) DtiLog_Price(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, d tmlDataElementing, price amibaModels.Price, err error) error {
	item := amibaModels.DtiModelingPrice{}
	item.EntId = ent.Id
	item.PurposeId = purpose.ID
	item.PeriodId = period.Id
	item.ItemCode = d.DataItemCode
	item.ItemId = d.DataItemId
	item.Date = d.DataDocDate
	item.ModelId = d.ModelingId
	item.FmGroupId = d.FmGroupId
	item.ToGroupId = d.ToGroupId
	item.Price = price.CostPrice
	if err != nil {
		item.Msg = err.Error()
	}
	if _, err := s.repo.Insert(&item); err != nil {
		glog.Errorf("企业:%v,核算:%v,期间:%v,更新取价日志错误,%s", ent.Name, purpose.Name, period.Name, err)
	}

	return nil
}
