package amibaServices

import (
	"time"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
)

func (s *modelSv) getBizData(m tmlModelingLine) []tmlDataElementing {
	if m.ModelLine.BizTypeEnum == "voucher" || m.ModelLine.BizTypeEnum == "" {
		return nil
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
	query.Where("d.ent_id=?", m.EntId).Where("d.doc_date between ? and ?", m.Period.FromDate.Format("2006-01-02"), m.Period.ToDate.Format("2006-01-02"))
	//阿米巴条件过滤
	if m.MatchGroup.Id != "" && len(m.MatchGroup.Datas) > 0 {
		switch m.MatchGroup.TypeEnum {
		case "org":
			query.In("d."+m.ModelLine.MatchDirectionEnum+"_org", m.MatchGroup.GetDataCodes())
		case "dept":
			query.In("d."+m.ModelLine.MatchDirectionEnum+"_dept", m.MatchGroup.GetDataCodes())
		case "work":
			query.In("d."+m.ModelLine.MatchDirectionEnum+"_work", m.MatchGroup.GetDataCodes())
		}
	}
	//模型条件过滤
	query.Where("d.biz_type=?", m.ModelLine.BizTypeEnum)
	if m.ModelLine.DocTypeCode != "" {
		query.Where("d.doc_type=?", m.ModelLine.DocTypeCode)
	}
	if m.ModelLine.ItemCategoryCode != "" {
		query.Where("d.item_category=?", m.ModelLine.ItemCategoryCode)
	}
	if m.ModelLine.ItemCode != "" {
		query.Where("d.item=?", m.ModelLine.ItemCode)
	}
	if m.ModelLine.TraderCode != "" {
		query.Where("d.trader=?", m.ModelLine.TraderCode)
	}
	if m.ModelLine.ProjectCode != "" {
		query.Where("d.project=?", m.ModelLine.ProjectCode)
	}
	if m.ModelLine.Factor1 != "" {
		query.Where("d.factor1=?", m.ModelLine.Factor1)
	}
	if m.ModelLine.Factor2 != "" {
		query.Where("d.factor2=?", m.ModelLine.Factor2)
	}
	if m.ModelLine.Factor3 != "" {
		query.Where("d.factor3=?", m.ModelLine.Factor3)
	}
	if m.ModelLine.Factor4 != "" {
		query.Where("d.factor4=?", m.ModelLine.Factor4)
	}
	if m.ModelLine.Factor5 != "" {
		query.Where("d.factor5=?", m.ModelLine.Factor5)
	}
	err := query.Find(&datas)
	if err != nil {
		glog.Printf("query error :%s", err)
		return nil
	}
	glog.Printf("查询业务数据:%v条,time:%v Seconds", len(datas), time.Now().Sub(fm_time).Seconds())

	fm_time = time.Now()
	tmlDatas := make([]tmlDataElementing, 0)
	for _, d := range datas {
		tml := tmlDataElementing{
			EntId: m.EntId, PeriodId: m.Period.Id, PurposeId: m.ModelLine.PurposeId,
			ModelingId: m.ModelLine.ModelId, ModelingLineId: m.ModelLine.Id, MatchDirectionEnum: m.ModelLine.MatchDirectionEnum, MatchGroupId: m.ModelLine.MatchGroupId,
			DefFmGroupId: m.ModelLine.GroupId, DefToGroupId: m.ModelLine.ToGroupId, ElementId: m.ModelLine.ElementId, BizTypeEnum: m.ModelLine.BizTypeEnum,
			ValueTypeEnum: m.ModelLine.ValueTypeEnum, Adjust: m.ModelLine.Adjust,
			DataId: d.Id, DataType: "biz",
			DataFmOrg: d.FmOrg, DataFmDept: d.FmDept, DataFmWork: d.FmWork, DataFmTeam: d.FmTeam, DataFmPerson: d.FmPerson,
			DataToOrg: d.ToOrg, DataToDept: d.ToDept, DataToWork: d.ToWork, DataToTeam: d.ToTeam, DataToPerson: d.ToPerson,
			DataTraderCode: d.Trader, DataItemCode: d.Item, DataItemCategoryCode: d.ItemCategory, DataProjectCode: d.Project, DataAccountCode: d.Account, DataUom: d.Uom, DataCurrency: d.Currency,
			DataQty: d.Qty, DataMoney: d.Money,
		}
		s.model_sv_handTmlData(&tml, m)
		if !tml.Deleted {
			tmlDatas = append(tmlDatas, tml)
		}
	}

	glog.Printf("处理业务数据,time:%v Seconds", time.Now().Sub(fm_time).Seconds())
	return tmlDatas
}
