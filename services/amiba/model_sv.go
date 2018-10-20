package amibaServices

import (
	"errors"
	"fmt"
	"time"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/services/cbo"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/repositories"
)

type ModelSv interface {
	Modeling(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelIds []string) (bool, error)
}
type modelSv struct {
	repo      *repositories.ModelRepo
	entSv     cboServices.EntSv
	itemsv    cboServices.ItemSv
	purposeSv PurposeSv
	groupSv   GroupSv
	pricelSv  PricelSv
}

func NewModelSv(repo *repositories.ModelRepo, entSv cboServices.EntSv, itemsv cboServices.ItemSv, purposeSv PurposeSv, groupSv GroupSv, pricelSv PricelSv) ModelSv {
	return &modelSv{repo: repo, entSv: entSv, itemsv: itemsv, purposeSv: purposeSv, groupSv: groupSv, pricelSv: pricelSv}
}
func (s *modelSv) Modeling(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelIds []string) (bool, error) {
	//获取期间数据
	fm_time = time.Now()
	//获取模型数据
	fm_time = time.Now()
	models, f := s.GetModels(ent.Id, purpose.Id, modelIds)
	if !f || len(models) == 0 {
		err := errors.New(fmt.Sprintf("企业:%v,核算:%v,期间:%v,找不到模型数据,%v", ent.Name, purpose.Name, period.Name, modelIds))
		glog.Printf("model data error :%s", err)
		return false, err
	}
	glog.Printf("企业:%v,核算:%v,期间:%v,取得模型数据=%v条,time:%v Seconds", ent.Name, purpose.Name, period.Name, len(models), time.Now().Sub(fm_time).Seconds())
	for _, m := range models {
		s.DtiLog_Init(ent, purpose, period, m.Id)

		//删除取价日志
		sql := "delete from `suite_amiba_dti_modeling_prices` where ent_id=? and purpose_id=? and period_id=? and model_id=?"
		if _, err := s.repo.Exec(sql, ent.Id, purpose.Id, period.Id, m.Id); err != nil {
			glog.Printf("企业:%v,核算:%v,期间:%v,删除取价日志错误:%s", ent.Name, purpose.Name, period.Name, err)
		}
	}

	for _, m := range models {
		//获取业务数据
		tmlDatas := make([]tmlDataElementing, 0)
		if m.Lines == nil || len(m.Lines) == 0 {
			continue
		}
		fm_time2 := time.Now()
		log := ""
		s.DtiLog_Begin(ent, purpose, period, m.Id)

		glog.Printf("企业:%v,核算:%v,模型:%v,开始业务数据和财务数据建模...", ent.Name, purpose.Name, m.Name)
		for _, ml := range m.Lines {
			if ml.MatchGroupId == "" {
				err := errors.New(fmt.Sprintf("企业:%v,核算:%v,模型:%v,行:%v,找不到匹配方阿米巴:%s", ent.Name, purpose.Name, m.Name, ml.Id, ml.MatchGroupId))
				glog.Printf("group data error :%s", err)
				s.DtiLog_Error(ent, purpose, period, m.Id, err)
				return false, err
			}

			if ml.MatchGroup.Datas == nil && len(ml.MatchGroup.Datas) == 0 {
				err := errors.New(fmt.Sprintf("企业:%v,核算:%v,模型:%v,匹配方巴:%v,必须是末级，且需要有明细构成", ent.Name, purpose.Name, m.Name, ml.MatchGroup.Name))
				glog.Printf("group data error :%s", err)
				s.DtiLog_Error(ent, purpose, period, m.Id, err)
				return false, err
			}

			fm_time = time.Now()
			tml, err := s.getBizData(ent, purpose, period, ml)
			if err != nil {
				glog.Error("获取业务数据建模错误：", err)
				s.DtiLog_Error(ent, purpose, period, m.Id, err)
				break
			}
			if tml != nil && len(tml) > 0 {
				tmlDatas = append(tmlDatas, tml...)
			}
			log = fmt.Sprintf("企业:%v,核算:%v,模型:%v,业务数据建模:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Name, len(tml), time.Now().Sub(fm_time).Seconds())
			glog.Printf(log)
			s.DtiLog_Log(ent, purpose, period, m.Id, log)
			fm_time = time.Now()
			tml, err = s.getFiData(ent, purpose, period, ml)
			if err != nil {
				glog.Error("获取财务数据建模错误：", err)
				s.DtiLog_Error(ent, purpose, period, m.Id, err)
				break
			}
			if tml != nil && len(tml) > 0 {
				tmlDatas = append(tmlDatas, tml...)
			}
			log = fmt.Sprintf("企业:%v,核算:%v,模型:%v,财务数据建模:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Name, len(tml), time.Now().Sub(fm_time).Seconds())
			glog.Printf(log)
			s.DtiLog_Log(ent, purpose, period, m.Id, log)
		}
		log = fmt.Sprintf("企业:%v,核算:%v,模型:%v,数据建模完成:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Name, len(tmlDatas), time.Now().Sub(fm_time2).Seconds())
		glog.Printf(log)
		s.DtiLog_Log(ent, purpose, period, m.Id, log)
		//保存单据
		s.Savedoc(ent, purpose, period, tmlDatas, m)
		log = fmt.Sprintf("企业:%v,核算:%v,模型:%v,保存单据完成,time:%v Seconds", ent.Name, purpose.Name, m.Name, time.Now().Sub(fm_time2).Seconds())
		glog.Printf(log)
		s.DtiLog_Log(ent, purpose, period, m.Id, log)

		s.DtiLog_Success(ent, purpose, period, m.Id)
	}
	//获取业务数据
	return true, nil
}
