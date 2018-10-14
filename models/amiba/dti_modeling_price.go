package amibaModels

import "time"

type DtiModelingPrice struct {
	Id        int64
	EntId     string    `xorm:"varchar(200) 'ent_id'"`
	PurposeId string    `xorm:"varchar(200) 'purpose_id'"`
	PeriodId  string    `xorm:"varchar(200) 'period_id'"`
	ModelId   string    `xorm:"varchar(200) 'model_id'"`
	Date      time.Time `xorm:"timestamp 'date'"`
	FmGroupId string    `xorm:"varchar(200) 'fm_group_id'"`
	ToGroupId string    `xorm:"varchar(200) 'to_group_id'"`
	ItemCode  string    `xorm:"varchar(200) 'item_code'"`
	ItemId    string    `xorm:"varchar(200) 'item_id'"`
	Price     float64   `xorm:"decimal(24,8) 'price'"`
	Msg       string    `xorm:"text 'msg'"`
	CreatedAt time.Time `xorm:"created 'created_at'"`
}

func (s *DtiModelingPrice) TableName() string {
	return "suite_amiba_dti_modeling_prices"
}
