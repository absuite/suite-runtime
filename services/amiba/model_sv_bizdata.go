package amibaServices

import (
	"log"
	"time"

	"github.com/absuite/suite-runtime/models/amiba"
)

func getBizData(m tmlModelingLine) []tmlDataElementing {
	if m.Model.BizTypeEnum == "voucher" || m.Model.BizTypeEnum == "" {
		return nil
	}
	fm_time = time.Now()
	datas := make([]amibaModels.DataBiz, 0)
	query := model_sv.repo.Select(` 
		d.doc_no,d.doc_date,d.biz_type,d.doc_type,
		d.fm_org,d.fm_dept,d.fm_work,d.fm_team,d.fm_wh,d.fm_person,
		d.to_org,d.to_dept,d.to_work,d.to_team,d.to_wh,d.to_person,
		d.trader,d.item,d.item_category,d.project,d.currency,d.uom,d.qty,d.price,d.money,d.tax,
		d.factor1,d.factor2,d.factor3,d.factor4,d.factor5
		`)
	query.Table("suite_amiba_doc_bizs").Alias("d")
	query.Where("d.ent_id=?", m.EntId).Where("d.doc_date between ? and ?", m.Period.FromDate.Format("2006-01-02"), m.Period.ToDate.Format("2006-01-02"))
	//阿米巴条件过滤
	if m.MatchGroup.Id != "" && len(m.MatchGroup.Datas) > 0 {
		switch m.MatchGroup.TypeEnum {
		case "org":
			query.In("d."+m.Model.MatchDirectionEnum+"_org", m.MatchGroup.GetDataCodes())
		case "dept":
			query.In("d."+m.Model.MatchDirectionEnum+"_dept", m.MatchGroup.GetDataCodes())
		case "work":
			query.In("d."+m.Model.MatchDirectionEnum+"_work", m.MatchGroup.GetDataCodes())
		}
	}
	//模型条件过滤
	query.Where("d.biz_type=?", m.Model.BizTypeEnum)
	if m.Model.DocTypeCode != "" {
		query.Where("d.doc_type=?", m.Model.DocTypeCode)
	}
	if m.Model.ItemCategoryCode != "" {
		query.Where("d.item_category=?", m.Model.ItemCategoryCode)
	}
	if m.Model.ItemCode != "" {
		query.Where("d.item=?", m.Model.ItemCode)
	}
	if m.Model.TraderCode != "" {
		query.Where("d.trader=?", m.Model.TraderCode)
	}
	if m.Model.ProjectCode != "" {
		query.Where("d.project=?", m.Model.ProjectCode)
	}
	if m.Model.Factor1 != "" {
		query.Where("d.factor1=?", m.Model.Factor1)
	}
	if m.Model.Factor2 != "" {
		query.Where("d.factor2=?", m.Model.Factor2)
	}
	if m.Model.Factor3 != "" {
		query.Where("d.factor3=?", m.Model.Factor3)
	}
	if m.Model.Factor4 != "" {
		query.Where("d.factor4=?", m.Model.Factor4)
	}
	if m.Model.Factor5 != "" {
		query.Where("d.factor5=?", m.Model.Factor5)
	}
	err := query.Find(&datas)
	if err != nil {
		log.Printf("query error :%s", err)
		return nil
	}
	log.Printf("查询业务数据:%v条,time:%v Seconds", len(datas), time.Now().Sub(fm_time).Seconds())

	fm_time = time.Now()
	tmlDatas := make([]tmlDataElementing, 0)
	for _, d := range datas {
		tml := tmlDataElementing{
			EntId: m.EntId, PeriodId: m.Period.Id, PurposeId: m.Model.PurposeId,
			ModelingId: m.Model.Id, ModelingLineId: m.Model.LineId, MatchDirectionEnum: m.Model.MatchDirectionEnum, MatchGroupId: m.Model.MatchGroupId,
			DefFmGroupId: m.Model.GroupId, DefToGroupId: m.Model.ToGroupId, ElementId: m.Model.ElementId, BizTypeEnum: m.Model.BizTypeEnum,
			ValueTypeEnum: m.Model.ValueTypeEnum, Adjust: m.Model.Adjust,
			DataId: d.Id, DataType: "biz",
			DataFmOrg: d.FmOrg, DataFmDept: d.FmDept, DataFmWork: d.FmWork, DataFmTeam: d.FmTeam, DataFmPerson: d.FmPerson,
			DataToOrg: d.ToOrg, DataToDept: d.ToDept, DataToWork: d.ToWork, DataToTeam: d.ToTeam, DataToPerson: d.ToPerson,
			DataTrader: d.Trader, DataItemCode: d.Item, DataItemCategory: d.ItemCategory, DataProject: d.Project, DataAccount: d.Account, DataUom: d.Uom, DataCurrency: d.Currency,
			DataQty: d.Qty, DataMoney: d.Money,
		}
		model_sv_handTmlData(&tml, m)
		if !tml.Deleted {
			tmlDatas = append(tmlDatas, tml)
		}
	}

	log.Printf("处理业务数据,time:%v Seconds", time.Now().Sub(fm_time).Seconds())
	return tmlDatas
}
