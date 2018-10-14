package amibaModels

import "time"

type DtiModeling struct {
	Id        string    `xorm:"pk varchar(200) 'id'"`
	EntId     string    `xorm:"varchar(200) 'ent_id'"`
	PurposeId string    `xorm:"varchar(200) 'purpose_id'"`
	PeriodId  string    `xorm:"varchar(200) 'period_id'"`
	ModelId   string    `xorm:"varchar(200) 'model_id'"`
	Succeed   int       `xorm:"int 'succeed'"`
	Status    int       `xorm:"int 'status'"`
	StartTime time.Time `xorm:"datetime 'start_time'"`
	EndTime   time.Time `xorm:"datetime 'end_time'"`
	Msg       string    `xorm:"text 'msg'"`
	CreatedAt time.Time `xorm:"created 'created_at'"`
}

func (s *DtiModeling) TableName() string {
	return "suite_amiba_dti_modelings"
}
