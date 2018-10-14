package amibaModels

import (
	"time"
)

type Price struct {
	Id        string `xorm:"varchar(200) 'id'"`
	EntId     string `xorm:"varchar(200) 'ent_id'"`
	Code      string `xorm:"varchar(200) 'code'"`
	Name      string `xorm:"varchar(200) 'name'"`
	PurposeId string `xorm:"varchar(200) 'purpose_id'"`
	PeriodId  string `xorm:"varchar(200) 'period_id'"`
	FmGroupId string `xorm:"varchar(200) 'fm_group_id'"`
	ToGroupId string `xorm:"varchar(200) 'to_group_id'"`

	ItemId    string  `xorm:"varchar(200) 'item_id'"`
	ItemCode  string  `xorm:"varchar(200) 'item_code'"`
	CostPrice float64 `xorm:"decimal 'cost_price'"`

	FmDate time.Time `xorm:"decimal 'fm_date'"`
	ToDate time.Time `xorm:"decimal 'to_date'"`

	CacheKey string
}
