package amibaModels

type Model struct {
	Id                 string `xorm:"varchar(200) 'id'"`
	LineId             string `xorm:"varchar(200) 'line_id'"`
	Code               string `xorm:"varchar(200) 'code'"`
	Name               string `xorm:"varchar(200) 'name'"`
	PurposeId          string `xorm:"varchar(200) 'purpose_id'"`
	GroupId            string `xorm:"varchar(200) 'group_id'"`
	ElementId          string `xorm:"varchar(200) 'element_id'"`
	MatchDirectionEnum string `xorm:"varchar(200) 'match_direction_enum'"`
	MatchGroupId       string `xorm:"varchar(200) 'match_group_id'"`
	BizTypeEnum        string `xorm:"varchar(200) 'biz_type_enum'"`
	DocTypeId          string `xorm:"varchar(200) 'doc_type_id'"`
	DocTypeCode        string `xorm:"varchar(200) 'doc_type_code'"`
	ItemCategoryId     string `xorm:"varchar(200) 'item_category_id'"`
	ItemCategoryCode   string `xorm:"varchar(200) 'item_category_code'"`
	ProjectCode        string `xorm:"varchar(200) 'project_code'"`
	AccountCode        string `xorm:"varchar(200) 'account_code'"`
	TraderId           string `xorm:"varchar(200) 'trader_id'"`
	TraderCode         string `xorm:"varchar(200) 'trader_code'"`
	ItemId             string `xorm:"varchar(200) 'item_id'"`
	ItemCode           string `xorm:"varchar(200) 'item_code'"`
	Factor1            string `xorm:"varchar(200) 'factor1'"`
	Factor2            string `xorm:"varchar(200) 'factor2'"`
	Factor3            string `xorm:"varchar(200) 'factor3'"`
	Factor4            string `xorm:"varchar(200) 'factor4'"`
	Factor5            string `xorm:"varchar(200) 'factor5'"`
	ValueTypeEnum      string `xorm:"varchar(200) 'value_type_enum'"`
	Adjust             string `xorm:"varchar(200) 'adjust'"`
	ToGroupId          string `xorm:"varchar(200) 'to_group_id'"`
	PriceId            string `xorm:"varchar(200) 'price_id'"`
}
