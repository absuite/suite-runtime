package amibaServices

import (
	"errors"
	"strconv"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/models/cbo"
)

func (s *modelSv) priceData(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, d *tmlDataElementing, m amibaModels.ModelLine) error {

	//依据阿米巴定义找来源阿米巴
	var fmGroup amibaModels.Group
	var toGroup amibaModels.Group
	if m.Group.TypeEnum == "org" && d.DataFmOrgCode != "" {
		fmGroup, _ = s.groupSv.FindByLineCode(ent.Id, m.PurposeId, d.DataFmOrgCode)
	}
	if m.Group.TypeEnum == "dept" && d.DataFmDeptCode != "" {
		fmGroup, _ = s.groupSv.FindByLineCode(ent.Id, m.PurposeId, d.DataFmDeptCode)
	}
	if m.Group.TypeEnum == "work" && d.DataFmWorkCode != "" {
		fmGroup, _ = s.groupSv.FindByLineCode(ent.Id, m.PurposeId, d.DataFmWorkCode)
	}
	if fmGroup.Id != "" {
		d.DataFmGroupId = fmGroup.Id
	}
	if m.Group.TypeEnum == "org" && d.DataToOrgCode != "" {
		toGroup, _ = s.groupSv.FindByLineCode(ent.Id, m.PurposeId, d.DataToOrgCode)
	}
	if m.Group.TypeEnum == "dept" && d.DataToDeptCode != "" {
		toGroup, _ = s.groupSv.FindByLineCode(ent.Id, m.PurposeId, d.DataToDeptCode)
	}
	if m.Group.TypeEnum == "work" && d.DataToWorkCode != "" {
		toGroup, _ = s.groupSv.FindByLineCode(ent.Id, m.PurposeId, d.DataToWorkCode)
	}
	if toGroup.Id != "" {
		d.DataToGroupId = fmGroup.Id
	}
	//数据巴的来源和去向相同时，则需要删除
	if d.DataToGroupId == d.DataFmGroupId {
		d.Deleted = true
		return nil
	}
	/*
		默认头上巴，为模型中指定的头上的巴,对方巴为模型中指定的对方巴
		匹配方：与原始业务数据中的巴进行匹配，并将匹配结果作为表头巴的建模成果。
		交易方：标识表头巴的交易对方巴是谁，当交易方为空时，使用原始业务数据匹配的巴，如果指定了，则直接使用交易巴。
		更新模型巴
		如果模型行定义了匹配方，则按匹配方更新巴
	*/
	if m.GroupId != "" {
		d.FmGroupId = m.GroupId
	} else if m.MatchDirectionEnum == "fm" {
		d.FmGroupId = d.DataFmGroupId
	} else if m.MatchDirectionEnum == "to" {
		d.FmGroupId = d.DataToGroupId
	}

	if m.ToGroupId != "" {
		d.ToGroupId = m.ToGroupId
	} else if m.MatchDirectionEnum == "fm" {
		d.ToGroupId = d.DataToGroupId
	} else if m.MatchDirectionEnum == "to" {
		d.ToGroupId = d.DataFmGroupId
	}
	/*来源和去向相同时，则需要删除，2、数据为空时，需要删除*/
	if d.FmGroupId == "" || d.FmGroupId == d.ToGroupId {
		d.Deleted = true
		return nil
	}
	//数量
	d.Qty = d.DataQty
	//金额
	if d.ValueTypeEnum == "qtyvalue" {
		d.Money = d.DataQty
	} else if d.ValueTypeEnum == "qty" {
		if d.Qty == 0 {
			d.Deleted = true
			glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,物料:%v,业务数据需要取价，但数量为0，被丢弃", ent.Name, purpose.Name, m.Id, period.Name, d.DataItemCode)
			return nil
		}
		/*如果取数量，则需要从价表里取单价*/
		price, f, err := s.pricelSv.GetPrice(PriceFind{EntId: d.EntId, PurposeId: d.PurposeId, PriceId: m.PriceId, FmGroupId: d.FmGroupId, ToGroupId: d.ToGroupId, ItemCode: d.DataItemCode, Date: d.DataDocDate})
		if err != nil {
			s.DtiLog_Price(ent, purpose, period, *d, price, err)
			d.Deleted = true
			return err
		}
		if !f {
			d.Deleted = true
			glog.Errorf("企业:%v,核算:%v,模型:%v,期间:%v,价表:%v,来源巴:%v,去向巴:%v,物料:%v,取不到价，find price:%s", ent.Name, purpose.Name, m.Id, period.Name, m.PriceId, d.FmGroupId, d.ToGroupId, d.DataItemCode, err)
			s.DtiLog_Price(ent, purpose, period, *d, price, errors.New("取不到价"))
			return nil
		}
		d.Price = price.CostPrice
		d.Money = d.Qty * price.CostPrice
	} else {
		d.Money = d.DataMoney
	}

	//计算调整比例
	if d.Adjust != "" && d.Adjust != "100" {
		justValue, err := strconv.Atoi(d.Adjust)
		if err != nil {
			d.Deleted = true
			glog.Error("企业:%v,核算:%v,模型:%v,期间:%v,物料:%v,模型调整错误，字符串转换成整数失败:%s，find price:%s", ent.Name, purpose.Name, m.Id, period.Name, d.DataItemCode, err)
			return err
		}
		d.Money = d.Money * float64(justValue) / 100
	}
	if d.Qty == 0 && d.Money == 0 {
		d.Deleted = true
		glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,物料:%v,数量和金额同时为0，数据被丢弃，find price:%s", ent.Name, purpose.Name, m.Id, period.Name, d.DataItemCode)
		return nil
	}
	return nil
}
