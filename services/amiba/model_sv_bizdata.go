package amibaServices

import (
	"time"

	"github.com/absuite/suite-runtime/models/cbo"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
)

func (s *modelSv) getBizData(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, m amibaModels.ModelLine) ([]tmlDataElementing, error) {
	if m.BizTypeEnum == "voucher" || m.BizTypeEnum == "" {
		return nil, nil
	}
	fm_time = time.Now()
	datas := make([]amibaModels.DataBiz, 0)
	query := s.repo.Select(` 
		max(d.doc_date) as doc_date,sum(d.qty) as qty,sum(d.money) as money,sum(d.tax) as tax,
		d.biz_type,d.doc_type,
		d.fm_org,d.fm_dept,d.fm_work,d.fm_team,d.fm_wh,d.fm_person,
		d.to_org,d.to_dept,d.to_work,d.to_team,d.to_wh,d.to_person,
		d.trader,d.item,d.item_category,d.project,d.currency,d.uom,
		d.factor1,d.factor2,d.factor3,d.factor4,d.factor5		
		`)
	query.Table("suite_amiba_doc_bizs").Alias("d")
	query.GroupBy(`
		d.biz_type,d.doc_type,
		d.fm_org,d.fm_dept,d.fm_work,d.fm_team,d.fm_wh,d.fm_person,
		d.to_org,d.to_dept,d.to_work,d.to_team,d.to_wh,d.to_person,
		d.trader,d.item,d.item_category,d.project,d.currency,d.uom,
		d.factor1,d.factor2,d.factor3,d.factor4,d.factor5
		`)
	query.Where("d.ent_id=?", ent.Id).Where("d.doc_date between ? and ?", period.FromDate.Format("2006-01-02"), period.ToDate.Format("2006-01-02"))
	//阿米巴条件过滤
	if m.MatchGroup.Id != "" && len(m.MatchGroup.Datas) > 0 {
		switch m.MatchGroup.TypeEnum {
		case "org":
			query.In("d."+m.MatchDirectionEnum+"_org", m.MatchGroup.GetDataCodes())
		case "dept":
			query.In("d."+m.MatchDirectionEnum+"_dept", m.MatchGroup.GetDataCodes())
		case "work":
			query.In("d."+m.MatchDirectionEnum+"_work", m.MatchGroup.GetDataCodes())
		}
	}
	//模型条件过滤
	query.Where("d.biz_type=?", m.BizTypeEnum)
	if m.DocTypeCode != "" {
		query.Where("d.doc_type=?", m.DocTypeCode)
	}
	if m.ItemCategoryCode != "" {
		query.Where("d.item_category=?", m.ItemCategoryCode)
	}
	if m.ItemCode != "" {
		query.Where("d.item=?", m.ItemCode)
	}
	if m.TraderCode != "" {
		query.Where("d.trader=?", m.TraderCode)
	}
	if m.ProjectCode != "" {
		query.Where("d.project=?", m.ProjectCode)
	}
	if m.Factor1 != "" {
		query.Where("d.factor1=?", m.Factor1)
	}
	if m.Factor2 != "" {
		query.Where("d.factor2=?", m.Factor2)
	}
	if m.Factor3 != "" {
		query.Where("d.factor3=?", m.Factor3)
	}
	if m.Factor4 != "" {
		query.Where("d.factor4=?", m.Factor4)
	}
	if m.Factor5 != "" {
		query.Where("d.factor5=?", m.Factor5)
	}
	err := query.Find(&datas)
	if err != nil {
		glog.Printf("query error :%s", err)
		return nil, err
	}
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,查询业务数据:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Id, period.Name, len(datas), time.Now().Sub(fm_time).Seconds())

	fm_time = time.Now()
	tmlDatas := make([]tmlDataElementing, 0)
	for _, d := range datas {
		tml := tmlDataElementing{
			EntId: ent.Id, PeriodId: period.Id, PurposeId: m.PurposeId,
			ModelingId: m.ModelId, ModelingLineId: m.Id, MatchDirectionEnum: m.MatchDirectionEnum, MatchGroupId: m.MatchGroupId,
			DefFmGroupId: m.GroupId, DefToGroupId: m.ToGroupId, ElementId: m.ElementId, BizTypeEnum: m.BizTypeEnum,
			ValueTypeEnum: m.ValueTypeEnum, Adjust: m.Adjust,
			DataId: d.Id, DataType: "biz",
			DataFmOrg: d.FmOrg, DataFmDept: d.FmDept, DataFmWork: d.FmWork, DataFmTeam: d.FmTeam, DataFmPerson: d.FmPerson,
			DataToOrg: d.ToOrg, DataToDept: d.ToDept, DataToWork: d.ToWork, DataToTeam: d.ToTeam, DataToPerson: d.ToPerson,
			DataTraderCode: d.Trader, DataItemCode: d.Item, DataItemCategoryCode: d.ItemCategory, DataProjectCode: d.Project, DataAccountCode: d.Account, DataUom: d.Uom, DataCurrency: d.Currency,
			DataQty: d.Qty, DataMoney: d.Money,
		}
		if err := s.model_sv_handTmlData(ent, purpose, period, &tml, m); err != nil {
			return nil, err
		}
		if !tml.Deleted {
			tmlDatas = append(tmlDatas, tml)
		}
	}
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,处理业务数据:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Id, period.Name, len(tmlDatas), time.Now().Sub(fm_time).Seconds())

	return tmlDatas, nil
}
