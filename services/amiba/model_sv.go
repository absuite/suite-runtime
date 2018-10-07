package amibaServices

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/repositories"
)

type ModelSv interface {
	Modeling(model amibaModels.Modeling) (bool, error)
}
type modelSv struct {
	repo *repositories.ModelRepo
}

var model_sv *modelSv

func NewModelSv(repo *repositories.ModelRepo) ModelSv {
	model_sv = &modelSv{repo: repo}
	return model_sv
}
func (s *modelSv) Modeling(model amibaModels.Modeling) (bool, error) {

	//获取期间数据
	fm_time = time.Now()
	period, f := model_sv_getPeriodData(model)
	if !f {
		err := errors.New(fmt.Sprintf("找不到期间数据:%s", model.PeriodId))
		log.Printf("period data error :%s", err)
		return false, err
	}
	log.Printf("获取期间数据,time:%v Seconds", time.Now().Sub(fm_time).Seconds())

	//获取模型数据
	fm_time = time.Now()
	modelLines, f := model_sv_getModelsData(model)
	if !f || len(modelLines) == 0 {
		err := errors.New(fmt.Sprintf("找不到模型数据:%s", model.ModelId))
		log.Printf("model data error :%s", err)
		return false, err
	}
	log.Printf("获取模型数据:%v条,time:%v Seconds", len(modelLines), time.Now().Sub(fm_time).Seconds())

	//获取阿米巴数据
	fm_time = time.Now()
	groups, f := model_sv_getGroups(model)
	if !f || len(groups) == 0 {
		err := errors.New(fmt.Sprintf("找不到阿米巴数据:%s", model.PurposeId))
		log.Printf("group data error :%s", err)
		return false, err
	}
	log.Printf("获取模型数据:%v条,time:%v Seconds", len(groups), time.Now().Sub(fm_time).Seconds())

	//获取业务数据
	tmlDatas := make([]tmlDataElementing, 0)
	for _, v := range modelLines {
		tmlModeling := tmlModelingLine{EntId: model.EntId, Period: period, Model: v, AllGroups: groups}

		group, found := model_sv_getGroup(v.GroupId, groups)
		if !found {
			err := errors.New(fmt.Sprintf("找不到阿米巴:%s", v.GroupId))
			log.Printf("group data error :%s", err)
			return false, err
		}
		tmlModeling.Group = group

		matchGroup, found := model_sv_getGroup(v.MatchGroupId, groups)
		if !found {
			err := errors.New(fmt.Sprintf("找不到匹配方阿米巴:%s", v.GroupId))
			log.Printf("group data error :%s", err)
			return false, err
		} else {
			if matchGroup.Datas == nil && len(matchGroup.Datas) == 0 {
				err := errors.New(fmt.Sprintf("匹配方巴必须是末级，且需要有明细构成:%s", matchGroup.Name))
				log.Printf("group data error :%s", err)
				return false, err
			}
			tmlModeling.MatchGroup = matchGroup
		}

		fm_time = time.Now()
		tml := getBizData(tmlModeling)
		if tml != nil && len(tml) > 0 {
			tmlDatas = append(tmlDatas, tml...)
		}
		log.Printf("业务数据建模:%v条,time:%v Seconds", len(tml), time.Now().Sub(fm_time).Seconds())

		fm_time = time.Now()
		tml = getFiData(tmlModeling)
		if tml != nil && len(tml) > 0 {
			tmlDatas = append(tmlDatas, tml...)
		}
		log.Printf("财务数据建模:%v条,time:%v Seconds", len(tml), time.Now().Sub(fm_time).Seconds())
	}
	//获取业务数据
	return true, nil
}
