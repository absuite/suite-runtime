package amibaModels

type Modeling struct {
	EntId     string `xorm:"varchar(200) 'ent_id'" json:"ent_id"`
	PurposeId string `xorm:"varchar(200) 'purpose_id'" json:"purpose_id"`
	PeriodId  string `xorm:"varchar(200) 'period_id'" json:"period_id"`
	ModelId   string `xorm:"varchar(200) 'model_ids'" json:"model_id"`
	Memo      string `xorm:"varchar(200) 'memo'" json:"memo"`
}
