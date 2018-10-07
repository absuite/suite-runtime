package amibaModels

type Modeling struct {
	EntId     string `xorm:"varchar(200) 'ent_id'"`
	PurposeId string `xorm:"varchar(200) 'purpose_id'"`
	PeriodId  string `xorm:"varchar(200) 'period_id'"`
	ModelId   string `xorm:"varchar(200) 'model_ids'"`
	Memo      string `xorm:"varchar(200) 'memo'"`
}
