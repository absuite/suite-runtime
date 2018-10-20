package amibaServices

import (
	"time"

	"github.com/absuite/suite-runtime/models/cbo"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
)

func (s *modelSv) getFiData(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, m amibaModels.ModelLine) ([]tmlDataElementing, error) {
	if m.BizTypeEnum != "voucher" || m.BizTypeEnum == "" {
		return nil, nil
	}
	fm_time = time.Now()
	datas := make([]amibaModels.DataBiz, 0)
	query := s.repo.Select(` 
		max(d.doc_date) as doc_date,sum(d.debit_money) as debit_money,sum(d.credit_money) as credit_money,
		d.biz_type,d.doc_type,
		d.fm_org,d.fm_dept,d.fm_work,d.fm_team,d.fm_wh,d.fm_person,
		d.trader,d.project,d.account,d.currency,
		d.factor1,d.factor2,d.factor3,d.factor4,d.factor5		
		`)
	query.Table("suite_amiba_doc_fis").Alias("d")
	query.GroupBy(`
		d.biz_type,d.doc_type,
		d.fm_org,d.fm_dept,d.fm_work,d.fm_team,d.fm_wh,d.fm_person,
		d.trader,d.project,d.account,d.currency,
		d.factor1,d.factor2,d.factor3,d.factor4,d.factor5		
		`)
	query.Where("d.ent_id=?", ent.Id).Where("d.doc_date between ? and ?", period.FromDate.Format("2006-01-02"), period.ToDate.Format("2006-01-02"))
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
	query.Where("d.biz_type=?", m.BizTypeEnum)
	if m.DocTypeCode != "" {
		query.Where("d.doc_type=?", m.DocTypeCode)
	}
	if m.AccountCode != "" {
		query.Where("d.account=?", m.AccountCode)
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
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,查询财务数据:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Id, period.Name, len(datas), time.Now().Sub(fm_time).Seconds())

	fm_time = time.Now()
	tmlDatas := make([]tmlDataElementing, 0)
	for _, d := range datas {
		tml := tmlDataElementing{
			EntId: ent.Id, PeriodId: period.Id, PurposeId: m.PurposeId,
			ModelingId: m.ModelId, ModelingLineId: m.Id, MatchDirectionEnum: m.MatchDirectionEnum, MatchGroupId: m.MatchGroupId,
			DefFmGroupId: m.GroupId, DefToGroupId: m.ToGroupId, ElementId: m.ElementId, BizTypeEnum: m.BizTypeEnum,
			ValueTypeEnum: m.ValueTypeEnum, Adjust: m.Adjust,
			DataId: d.Id, DataType: "fi",
			DataTraderCode: d.Trader,
			DataItemCode:   d.Item, DataItemCategoryCode: d.ItemCategory,
			DataProjectCode: d.Project, DataAccountCode: d.Account, DataUomCode: d.Uom, DataCurrencyCode: d.Currency,
			DataQty: d.Qty,
		}
		if v, f := s.itemsv.FindByCode(ent.Id, tml.DataItemCode); f {
			tml.DataItemId = v.Id
		}
		if v, f := s.traderSv.FindByCode(ent.Id, tml.DataTraderCode); f {
			tml.DataTraderId = v.Id
		}
		if v, f := s.projectSv.FindByCode(ent.Id, tml.DataProjectCode); f {
			tml.DataProjectId = v.Id
		}
		if m.ValueTypeEnum == "debit" {
			tml.DataMoney = d.DebitMoney
		} else {
			tml.DataMoney = d.CreditMoney
		}
		if m.MatchDirectionEnum == "fm" {
			tml.DataFmOrgCode = d.FmOrg
			tml.DataFmDeptCode = d.FmDept
			tml.DataFmWorkCode = d.FmWork
			tml.DataFmTeamCode = d.FmTeam
			tml.DataFmPersonCode = d.FmPerson
		} else {
			tml.DataToOrgCode = d.ToOrg
			tml.DataToDeptCode = d.ToDept
			tml.DataToWorkCode = d.ToWork
			tml.DataToTeamCode = d.ToTeam
			tml.DataToPersonCode = d.ToPerson
		}
		if err := s.priceData(ent, purpose, period, &tml, m); err != nil {
			return nil, err
		}
		if !tml.Deleted {
			tmlDatas = append(tmlDatas, tml)
		}
	}
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,处理财务数据:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Id, period.Name, len(tmlDatas), time.Now().Sub(fm_time).Seconds())

	return tmlDatas, nil
}
