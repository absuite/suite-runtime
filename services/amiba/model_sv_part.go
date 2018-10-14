package amibaServices

import (
	"errors"
	"strconv"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/models/cbo"
)

func (s *modelSv) model_sv_handTmlData(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, d *tmlDataElementing, m amibaModels.ModelLine) error {

	//依据阿米巴定义找来源阿米巴
	var fmGroup amibaModels.Group
	var toGroup amibaModels.Group
	if m.Group.TypeEnum == "org" && d.DataFmOrg != "" {
		fmGroup, _ = s.GetGroupByLineCode(ent.Id, m.PurposeId, d.DataFmOrg)
	}
	if m.Group.TypeEnum == "dept" && d.DataFmDept != "" {
		fmGroup, _ = s.GetGroupByLineCode(ent.Id, m.PurposeId, d.DataFmDept)
	}
	if m.Group.TypeEnum == "work" && d.DataFmWork != "" {
		fmGroup, _ = s.GetGroupByLineCode(ent.Id, m.PurposeId, d.DataFmWork)
	}
	if fmGroup.Id != "" {
		d.DataFmGroupId = fmGroup.Id
	}
	if m.Group.TypeEnum == "org" && d.DataToOrg != "" {
		toGroup, _ = s.GetGroupByLineCode(ent.Id, m.PurposeId, d.DataToOrg)
	}
	if m.Group.TypeEnum == "dept" && d.DataToDept != "" {
		toGroup, _ = s.GetGroupByLineCode(ent.Id, m.PurposeId, d.DataToDept)
	}
	if m.Group.TypeEnum == "work" && d.DataToWork != "" {
		toGroup, _ = s.GetGroupByLineCode(ent.Id, m.PurposeId, d.DataToWork)
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
		price, f, err := price_sv.GetPrice(PriceFind{EntId: d.EntId, PurposeId: d.PurposeId, PriceId: m.PriceId, FmGroupId: d.FmGroupId, ToGroupId: d.ToGroupId, ItemCode: d.DataItemCode, Date: d.DataDocDate})
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
func (s *modelSv) CachePeriods(entId string) error {
	if model_sv_cache_periods[entId] == nil {
		model_sv_cache_periods[entId] = make(map[string]cboModels.Period)
	}
	periods := make([]cboModels.Period, 0)
	query := s.repo.Select("c.type_enum,c.id as calendar_id,p.id,p.code,p.name,p.year,p.from_date,p.to_date").Table("suite_cbo_period_calendars").Alias("c")
	query.Join("inner", []string{"suite_cbo_period_accounts", "p"}, "c.id=p.calendar_id")
	query.Where("p.ent_id = ? ", entId)
	if err := query.Find(&periods); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	for i, item := range periods {
		model_sv_cache_periods[entId][item.Id] = periods[i]
	}
	return nil
}

/**
*获取核算
 */
func (s *modelSv) GetEnt(entId string) (cboModels.Ent, bool) {
	m := cboModels.Ent{}
	query := s.repo.Select("c.id,c.code,c.name").Table("gmf_sys_ents").Alias("c")
	query.Where("c.id=?", entId)
	has, err := query.Get(&m)
	if err != nil {
		glog.Printf("query error :%s", err)
		return cboModels.Ent{}, false
	}
	return m, has
}

/**
*获取核算
 */
func (s *modelSv) GetPurpose(entId, purposeId string) (amibaModels.Purpose, bool) {
	m := amibaModels.Purpose{}
	query := s.repo.Select("c.id,c.code,c.name").Table("suite_amiba_purposes").Alias("c")
	query.Where("c.ent_id = ? and c.id=?", entId, purposeId)
	has, err := query.Get(&m)
	if err != nil {
		glog.Printf("query error :%s", err)
		return amibaModels.Purpose{}, false
	}
	return m, has
}

/**
* @api 获取期间集合
 */
func (s *modelSv) GetPeriod(entId, periodId string) (cboModels.Period, bool) {
	_, ok := model_sv_cache_periods[entId]
	if !ok {
		s.CachePeriods(entId)
		if _, ok := model_sv_cache_periods[entId]; !ok {
			return cboModels.Period{}, false
		}
	}
	if v, ok := model_sv_cache_periods[entId][periodId]; ok {
		return v, ok
	}
	return cboModels.Period{}, false
}

func (s *modelSv) CacheGroups(entId, purposeId string) error {
	if model_sv_cache_groups[entId] == nil {
		model_sv_cache_groups[entId] = make(map[string]map[string]amibaModels.Group)
	}
	model_sv_cache_groups[entId][purposeId] = make(map[string]amibaModels.Group)

	groups := make([]amibaModels.Group, 0)
	query := s.repo.Select("g.id,g.code,g.name,g.type_enum").Table("suite_amiba_groups").Alias("g")
	query.Where("g.purpose_id = ? and g.ent_id = ? ", purposeId, entId)
	if err := query.Find(&groups); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if len(groups) == 0 {
		return nil
	}
	groupData := make([]amibaModels.GroupData, 0)
	groupDataItem := make([]amibaModels.GroupData, 0)

	query = s.repo.Select("g.id as group_id,g.type_enum as type_enum,d.id,d.code,d.name").Table("suite_amiba_groups").Alias("g")
	query.Join("inner", []string{"suite_amiba_group_lines", "gl"}, "g.id=gl.group_id")
	query.Join("inner", []string{"suite_cbo_orgs", "d"}, "gl.data_id=d.id").Where("g.type_enum=?", "org")
	if err := query.Find(&groupDataItem); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if len(groupDataItem) > 0 {
		groupData = append(groupData, groupDataItem...)
	}

	query = s.repo.Select("g.id as group_id,g.type_enum as type_enum,d.id,d.code,d.name").Table("suite_amiba_groups").Alias("g")
	query.Join("inner", []string{"suite_amiba_group_lines", "gl"}, "g.id=gl.group_id")
	query.Join("inner", []string{"suite_cbo_depts", "d"}, "gl.data_id=d.id").Where("g.type_enum=?", "dept")
	if err := query.Find(&groupDataItem); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if len(groupDataItem) > 0 {
		groupData = append(groupData, groupDataItem...)
	}

	query = s.repo.Select("g.id as group_id,g.type_enum as type_enum,d.id,d.code,d.name").Table("suite_amiba_groups").Alias("g")
	query.Join("inner", []string{"suite_amiba_group_lines", "gl"}, "g.id=gl.group_id")
	query.Join("inner", []string{"suite_cbo_works", "d"}, "gl.data_id=d.id").Where("g.type_enum=?", "work")
	if err := query.Find(&groupDataItem); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if len(groupDataItem) > 0 {
		groupData = append(groupData, groupDataItem...)
	}
	if len(groupData) > 0 {
		for gi, gv := range groups {
			for _, vv := range groupData {
				if gv.Id == vv.GroupId {
					groups[gi].AddData(vv)
				}
			}
		}
	}
	for i, group := range groups {
		model_sv_cache_groups[entId][purposeId][group.Id] = groups[i]
	}
	return nil
}

/**
* @api 获取阿米巴集合
 */
func (s *modelSv) GetGroups(entId, purposeId string) ([]amibaModels.Group, bool) {
	if model_sv_cache_groups[entId] == nil || model_sv_cache_groups[entId][purposeId] == nil {
		s.CacheGroups(entId, purposeId)
		if model_sv_cache_groups[entId] == nil || model_sv_cache_groups[entId][purposeId] == nil {
			return nil, false
		}
	}
	items := make([]amibaModels.Group, 0)
	maps := model_sv_cache_groups[entId][purposeId]
	if len(maps) > 0 {
		for _, m := range maps {
			items = append(items, m)
		}
	}
	return items, len(items) > 0
}

/**
* @api 获取阿米巴
 */
func (s *modelSv) GetGroup(entId, purposeId, groupId string) (amibaModels.Group, bool) {
	if groupId == "" || entId == "" || purposeId == "" {
		return amibaModels.Group{}, false
	}
	if model_sv_cache_groups[entId] == nil || model_sv_cache_groups[entId][purposeId] == nil {
		s.CacheGroups(entId, purposeId)
		if model_sv_cache_groups[entId] == nil || model_sv_cache_groups[entId][purposeId] == nil {
			return amibaModels.Group{}, false
		}
	}
	return model_sv_cache_groups[entId][purposeId][groupId], true
}

/**
* @api 通过数据组成获取阿米巴集合
 */
func (s *modelSv) GetGroupByLineCode(entId, purposeId, code string) (amibaModels.Group, bool) {
	if code == "" || entId == "" || purposeId == "" {
		return amibaModels.Group{}, false
	}
	groups, f := s.GetGroups(entId, purposeId)
	if !f {
		return amibaModels.Group{}, false
	}
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
	return amibaModels.Group{}, false
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
			if gv, f := s.GetGroup(entId, purposeId, item.GroupId); f {
				results[i].Group = gv
			}
			for _, lv := range resultLines {
				if lv.ModelId == item.Id {
					if gv, f := s.GetGroup(entId, purposeId, lv.GroupId); f {
						lv.Group = gv
					}
					if gv, f := s.GetGroup(entId, purposeId, lv.ToGroupId); f {
						lv.ToGroup = gv
					}
					if gv, f := s.GetGroup(entId, purposeId, lv.MatchGroupId); f {
						lv.MatchGroup = gv
					}
					results[i].AddLine(lv)
				}
			}
		}
	}
	return results, true
}
