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

type tmlDataElementing struct {
	EntId     string
	PurposeId string
	PeriodId  string
	FmGroupId string //结果来源巴
	ToGroupId string //结果目标巴

	ModelingId         string //模型行ID
	ModelingLineId     string //模型行ID
	MatchDirectionEnum string
	MatchGroupId       string //匹配方：与原始业务数据中的巴进行匹配，并将匹配结果作为表头巴的建模成果。
	DefFmGroupId       string //模型头定义的巴
	DefToGroupId       string //模型头定义的巴,交易方：标识表头巴的交易对方巴是谁，当交易方为空时，使用原始业务数据匹配的巴，如果指定了，则直接使用交易巴。

	ElementId     string
	BizTypeEnum   string //业务类型
	ValueTypeEnum string
	Adjust        string

	DataId        string
	DataType      string
	DataFmGroupId string //业务数据来源对应的巴
	DataToGroupId string //业务数据目标对应的巴

	DataDocNo        string
	DataDocDate      time.Time
	DataFmOrgCode    string
	DataFmDeptCode   string
	DataFmWorkCode   string
	DataFmTeamCode   string
	DataFmWhCode     string
	DataFmPersonCode string

	DataToOrgCode    string
	DataToDeptCode   string
	DataToWorkCode   string
	DataToTeamCode   string
	DataToWhCode     string
	DataToPersonCode string

	DataTraderId         string
	DataTraderCode       string
	DataItemCode         string
	DataItemId           string
	DataItemCategoryCode string
	DataItemCategoryId   string
	DataProjectId        string
	DataProjectCode      string
	DataAccountCode      string
	DataCurrencyCode     string
	DataUomCode          string
	DataQty              float64
	DataMoney            float64

	Qty   float64
	Price float64
	Money float64

	Deleted bool
}

type ModelSv interface {
	Modeling(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, modelIds []string) (bool, error)
}
type modelSv struct {
	repo      *repositories.ModelRepo
	entSv     cboServices.EntSv
	itemsv    cboServices.ItemSv
	traderSv  cboServices.TraderSv
	projectSv cboServices.ProjectSv
	purposeSv PurposeSv
	groupSv   GroupSv
	pricelSv  PricelSv
}

func NewModelSv(repo *repositories.ModelRepo, entSv cboServices.EntSv, itemsv cboServices.ItemSv, traderSv cboServices.TraderSv, projectSv cboServices.ProjectSv, purposeSv PurposeSv, groupSv GroupSv, pricelSv PricelSv) ModelSv {
	return &modelSv{repo: repo, entSv: entSv, itemsv: itemsv, purposeSv: purposeSv, groupSv: groupSv, pricelSv: pricelSv, projectSv: projectSv, traderSv: traderSv}
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

/**
* @api 获取模型集合
 */
func (s *modelSv) GetModels(entId, purposeId string, modelIds []string) ([]amibaModels.Model, bool) {
	results := make([]amibaModels.Model, 0)
	query := s.repo.Select(`m.id,m.code,m.name,m.purpose_id,m.group_id`).Table("suite_amiba_modelings").Alias("m")
	if modelIds != nil && len(modelIds) > 0 {
		query.In("m.id", modelIds)
	}
	if purposeId != "" {
		query.Where("m.purpose_id=?", purposeId)
	}
	if entId != "" {
		query.Where("m.ent_id=?", entId)
	}
	if err := query.Find(&results); err != nil {
		glog.Printf("query error :%s", err)
		return nil, false
	}
	if len(results) == 0 {
		glog.Printf("查询不到模型数据:ent_id=%v,purpose_id=%v,modelIds=%v", entId, purposeId, modelIds)
		return results, false
	}
	resultLines := make([]amibaModels.ModelLine, 0)
	query = s.repo.Select(`m.id as model_id,m.purpose_id,m.group_id,
		ml.id as id,ml.element_id,ml.match_direction_enum,ml.match_group_id,
		ml.biz_type_enum,ml.doc_type_id,ml.item_category_id,itemc.code as item_category_code,
		ml.account_code,ml.project_code,
		ml.trader_id,trader.code as trader_code,ml.item_id,item.code as item_code,ml.factor1,ml.factor2,ml.factor3,ml.factor4,ml.factor5,
		ml.value_type_enum,ml.adjust,
		ml.to_group_id,ml.price_id`).Table("suite_amiba_modelings").Alias("m")
	query.Join("inner", []string{"suite_amiba_modeling_lines", "ml"}, "m.id=ml.modeling_id")
	query.Join("left", []string{"suite_cbo_items", "item"}, "ml.item_id=item.id")
	query.Join("left", []string{"suite_cbo_item_categories", "itemc"}, "itemc.id=ml.item_category_id")
	query.Join("left", []string{"suite_cbo_traders", "trader"}, "ml.trader_id=trader.id")
	if modelIds != nil && len(modelIds) > 0 {
		query.In("m.id", modelIds)
	}
	if purposeId != "" {
		query.Where("m.purpose_id=?", purposeId)
	}
	if entId != "" {
		query.Where("m.ent_id=?", entId)
	}
	if err := query.Find(&resultLines); err != nil {
		glog.Printf("query error :%s", err)
		return nil, false
	}
	if len(resultLines) > 0 {
		for i, item := range results {
			if gv, f := s.groupSv.Get(entId, item.GroupId); f {
				results[i].Group = gv
			}
			for _, lv := range resultLines {
				if lv.ModelId == item.Id {
					if gv, f := s.groupSv.Get(entId, lv.GroupId); f {
						lv.Group = gv
					}
					if gv, f := s.groupSv.Get(entId, lv.ToGroupId); f {
						lv.ToGroup = gv
					}
					if gv, f := s.groupSv.Get(entId, lv.MatchGroupId); f {
						lv.MatchGroup = gv
					}
					results[i].AddLine(lv)
				}
			}
		}
	}
	return results, true
}
