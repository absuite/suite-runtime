package amibaServices

import (
	"strconv"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/models/cbo"
)

func (s *modelSv) model_sv_handTmlData(d *tmlDataElementing, c tmlModelingLine) {

	//依据阿米巴定义找来源阿米巴
	var fmGroup amibaModels.Group
	var toGroup amibaModels.Group
	if c.Group.TypeEnum == "org" && d.DataFmOrg != "" {
		fmGroup, _ = s.model_sv_getGroupByLineCode(d.DataFmOrg, c.AllGroups)
	}
	if c.Group.TypeEnum == "dept" && d.DataFmDept != "" {
		fmGroup, _ = s.model_sv_getGroupByLineCode(d.DataFmDept, c.AllGroups)
	}
	if c.Group.TypeEnum == "work" && d.DataFmWork != "" {
		fmGroup, _ = s.model_sv_getGroupByLineCode(d.DataFmWork, c.AllGroups)
	}
	if fmGroup.Id != "" {
		d.DataFmGroupId = fmGroup.Id
	}
	if c.Group.TypeEnum == "org" && d.DataToOrg != "" {
		toGroup, _ = s.model_sv_getGroupByLineCode(d.DataToOrg, c.AllGroups)
	}
	if c.Group.TypeEnum == "dept" && d.DataToDept != "" {
		toGroup, _ = s.model_sv_getGroupByLineCode(d.DataToDept, c.AllGroups)
	}
	if c.Group.TypeEnum == "work" && d.DataToWork != "" {
		toGroup, _ = s.model_sv_getGroupByLineCode(d.DataToWork, c.AllGroups)
	}
	if toGroup.Id != "" {
		d.DataToGroupId = fmGroup.Id
	}
	//数据巴的来源和去向相同时，则需要删除
	if d.DataToGroupId == d.DataFmGroupId {
		d.Deleted = true
		return
	}
	/*
		默认头上巴，为模型中指定的头上的巴,对方巴为模型中指定的对方巴
		匹配方：与原始业务数据中的巴进行匹配，并将匹配结果作为表头巴的建模成果。
		交易方：标识表头巴的交易对方巴是谁，当交易方为空时，使用原始业务数据匹配的巴，如果指定了，则直接使用交易巴。
		更新模型巴
		如果模型行定义了匹配方，则按匹配方更新巴
	*/
	if c.Model.GroupId != "" {
		d.FmGroupId = c.Model.GroupId
	} else if c.Model.MatchDirectionEnum == "fm" {
		d.FmGroupId = d.DataFmGroupId
	} else if c.Model.MatchDirectionEnum == "to" {
		d.FmGroupId = d.DataToGroupId
	}

	if c.Model.ToGroupId != "" {
		d.ToGroupId = c.Model.ToGroupId
	} else if c.Model.MatchDirectionEnum == "fm" {
		d.ToGroupId = d.DataToGroupId
	} else if c.Model.MatchDirectionEnum == "to" {
		d.ToGroupId = d.DataFmGroupId
	}
	/*来源和去向相同时，则需要删除，2、数据为空时，需要删除*/
	if d.FmGroupId == "" || d.FmGroupId == d.ToGroupId {
		d.Deleted = true
		return
	}
	//数量
	d.Qty = d.DataQty
	//金额
	if d.ValueTypeEnum == "qtyvalue" {
		d.Money = d.DataQty
	} else if d.ValueTypeEnum == "qty" {
		if d.Qty == 0 {
			d.Deleted = true
			glog.Printf("业务数据需要取价，将数量为0，被丢弃!")
			return
		}
		/*如果取数量，则需要从价表里取单价*/
		price, err := price_sv.GetPrice(PriceFind{EntId: d.EntId, PurposeId: d.PurposeId, FmGroupId: d.FmGroupId, ToGroupId: d.ToGroupId, ItemCode: d.DataItemCode, Date: d.DataDocDate})
		if err != nil {
			d.Deleted = true
			glog.Printf("find price:%s", err)
			return
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
			glog.Printf("模型调整错误，字符串转换成整数失败:%s", err)
			return
		}
		d.Money = d.Money * float64(justValue)
	}
	if d.Qty == 0 && d.Money == 0 {
		d.Deleted = true
		glog.Printf("数量和金额同时为0，数据被丢弃!")
		return
	}
}

func (s *modelSv) model_sv_getPeriodData(model amibaModels.Modeling) (m cboModels.Period, found bool) {
	query := s.repo.Select("c.type_enum,c.id as calendar_id,p.id,p.code,p.name,p.year,p.from_date,p.to_date").Table("suite_cbo_period_calendars").Alias("c")
	query.Join("inner", []string{"suite_cbo_period_accounts", "p"}, "c.id=p.calendar_id")
	query.Where("p.id = ? ", model.PeriodId)
	if _, err := query.Get(&m); err != nil {
		glog.Printf("query error :%s", err)
		found = false
	}
	if m.Id != "" {
		found = true
	}
	return
}
func (s *modelSv) model_sv_getGroupByLineCode(code string, groups []amibaModels.Group) (m amibaModels.Group, found bool) {
	for _, g := range groups {
		if g.Datas == nil || len(g.Datas) == 0 {
			continue
		}
		for _, d := range g.Datas {
			if d.Code == code {
				return g, true
			}
		}
	}
	return
}
func (s *modelSv) model_sv_getGroup(groupId string, groups []amibaModels.Group) (m amibaModels.Group, found bool) {
	for _, g := range groups {
		if g.Id == groupId {
			return g, true
		}
	}
	return
}
func (s *modelSv) model_sv_getGroups(model amibaModels.Modeling) (m []amibaModels.Group, found bool) {
	query := s.repo.Select("g.id,g.code,g.name,g.type_enum").Table("suite_amiba_groups").Alias("g")
	query.Where("g.purpose_id = ? and g.ent_id = ? ", model.PurposeId, model.EntId)
	err := query.Find(&m)
	if err != nil {
		glog.Printf("query error :%s", err)
		found = false
		return
	}
	if len(m) == 0 {
		found = false
		return
	}
	groupType := "org"
	for _, g := range m {
		if g.TypeEnum != "" {
			groupType = g.TypeEnum
			break
		}
	}
	found = true
	groupData := make([]amibaModels.GroupData, 0)
	query = s.repo.Select("g.id as group_id,g.type_enum as type_enum,d.id,d.code,d.name").Table("suite_amiba_groups").Alias("g")
	query.Join("inner", []string{"suite_amiba_group_lines", "gl"}, "g.id=gl.group_id")

	switch groupType {
	case "org":
		query.Join("inner", []string{"suite_cbo_orgs", "d"}, "gl.data_id=d.id")
	case "dept":
		query.Join("inner", []string{"suite_cbo_depts", "d"}, "gl.data_id=d.id")
	case "work":
		query.Join("inner", []string{"suite_cbo_works", "d"}, "gl.data_id=d.id")
	}
	err = query.Find(&groupData)
	if err != nil {
		glog.Printf("query error :%s", err)
		found = false
		return
	}
	if len(groupData) > 0 {
		for gi, gv := range m {
			for vi, vv := range groupData {
				if gv.Id == vv.GroupId {
					m[gi].AddData(groupData[vi])
					break
				}
			}
		}
	}
	return
}
func (s *modelSv) model_sv_getModelsData(model amibaModels.Modeling) (m []amibaModels.Model, found bool) {
	query := s.repo.Select(`m.id,m.code,m.name,m.purpose_id,m.group_id,
		ml.id as line_id,ml.element_id,ml.match_direction_enum,ml.match_group_id,
		ml.biz_type_enum,ml.doc_type_id,ml.item_category_id,itemc.code as item_category_code,
		ml.account_code,ml.project_code,
		ml.trader_id,trader.code as trader_code,ml.item_id,item.code as item_code,ml.factor1,ml.factor2,ml.factor3,ml.factor4,ml.factor5,
		ml.value_type_enum,ml.adjust,
		ml.to_group_id,ml.price_id`).Table("suite_amiba_modelings").Alias("m")
	query.Join("inner", []string{"suite_amiba_modeling_lines", "ml"}, "m.id=ml.modeling_id")
	query.Join("left", []string{"suite_cbo_items", "item"}, "ml.item_id=item.id")
	query.Join("left", []string{"suite_cbo_item_categories", "itemc"}, "itemc.id=ml.item_category_id")
	query.Join("left", []string{"suite_cbo_traders", "trader"}, "ml.trader_id=trader.id")
	query.Where("m.id=?", model.ModelId)
	if err := query.Find(&m); err != nil {
		glog.Printf("query error :%s", err)
		found = false
	}
	if len(m) > 0 {
		found = true
	}
	return
}
