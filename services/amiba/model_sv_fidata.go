package amibaServices

import (
	"log"
	"time"

	"github.com/absuite/suite-runtime/models/amiba"
)

func (s *modelSv) getFiData(m tmlModelingLine) []tmlDataElementing {
	if m.Model.BizTypeEnum != "voucher" || m.Model.BizTypeEnum == "" {
		return nil
	}
	fm_time = time.Now()
	datas := make([]amibaModels.DataBiz, 0)
	query := s.repo.Select(` 
		d.doc_no,d.doc_date,d.biz_type,d.doc_type,
		d.fm_org,d.fm_dept,d.fm_work,d.fm_team,d.fm_wh,d.fm_person,
		d.trader,d.project,d.account,d.currency,d.debit_money,d.credit_money,
		d.factor1,d.factor2,d.factor3,d.factor4,d.factor5
		`)
	query.Table("suite_amiba_doc_fis").Alias("d")
	query.Where("d.ent_id=?", m.EntId).Where("d.doc_date between ? and ?", m.Period.FromDate.Format("2006-01-02"), m.Period.ToDate.Format("2006-01-02"))
	//阿米巴条件过滤
	if m.MatchGroup.Id != "" && len(m.MatchGroup.Datas) > 0 {
		switch m.MatchGroup.TypeEnum {
		case "org":
			query.In("d.fm_org", m.MatchGroup.GetDataCodes())
		case "dept":
			query.In("d.fm_dept", m.MatchGroup.GetDataCodes())
		case "work":
			query.In("d.fm_work", m.MatchGroup.GetDataCodes())
		}
	}
	//模型条件过滤
	query.Where("d.biz_type=?", m.Model.BizTypeEnum)
	if m.Model.DocTypeCode != "" {
		query.Where("d.doc_type=?", m.Model.DocTypeCode)
	}
	if m.Model.AccountCode != "" {
		query.Where("d.account=?", m.Model.AccountCode)
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

	log.Printf("查询财务数据:%v条,time:%v Seconds", len(datas), time.Now().Sub(fm_time).Seconds())

	fm_time = time.Now()
	tmlDatas := make([]tmlDataElementing, 0)
	for _, d := range datas {
		tml := tmlDataElementing{
			EntId: m.EntId, PeriodId: m.Period.Id, PurposeId: m.Model.PurposeId,
			ModelingId: m.Model.Id, ModelingLineId: m.Model.LineId, MatchDirectionEnum: m.Model.MatchDirectionEnum, MatchGroupId: m.Model.MatchGroupId,
			DefFmGroupId: m.Model.GroupId, DefToGroupId: m.Model.ToGroupId, ElementId: m.Model.ElementId, BizTypeEnum: m.Model.BizTypeEnum,
			ValueTypeEnum: m.Model.ValueTypeEnum, Adjust: m.Model.Adjust,
			DataId: d.Id, DataType: "fi",
			DataTraderCode: d.Trader, DataItemCode: d.Item, DataItemCategoryCode: d.ItemCategory, DataProjectCode: d.Project, DataAccountCode: d.Account, DataUom: d.Uom, DataCurrency: d.Currency,
			DataQty: d.Qty,
		}
		if m.Model.ValueTypeEnum == "debit" {
			tml.DataMoney = d.DebitMoney
		} else {
			tml.DataMoney = d.CreditMoney
		}
		if m.Model.MatchDirectionEnum == "fm" {
			tml.DataFmOrg = d.FmOrg
			tml.DataFmDept = d.FmDept
			tml.DataFmWork = d.FmWork
			tml.DataFmTeam = d.FmTeam
			tml.DataFmPerson = d.FmPerson
		} else {
			tml.DataToOrg = d.ToOrg
			tml.DataToDept = d.ToDept
			tml.DataToWork = d.ToWork
			tml.DataToTeam = d.ToTeam
			tml.DataToPerson = d.ToPerson
		}
		s.model_sv_handTmlData(&tml, m)
		if !tml.Deleted {
			tmlDatas = append(tmlDatas, tml)
		}
	}

	log.Printf("处理财务数据,time:%v Seconds", time.Now().Sub(fm_time).Seconds())
	return tmlDatas
}
