package amibaModels

import "time"

type DataDoc struct {
	Id    string `xorm:"pk varchar(200) 'id'"`
	EntId string `xorm:"varchar(200) 'ent_id'"`

	DocNo       string    `xorm:"varchar(200) 'doc_no'"`
	DocDate     time.Time `xorm:"timestamp 'doc_date'"`
	PurposeId   string    `xorm:"varchar(200) 'purpose_id'"`
	PeriodId    string    `xorm:"varchar(200) 'period_id'"`
	UseType     string    `xorm:"varchar(200) 'use_type_enum'"`
	FmGroupId   string    `xorm:"varchar(200) 'fm_group_id'"`
	ToGroupId   string    `xorm:"varchar(200) 'to_group_id'"`
	ElementId   string    `xorm:"varchar(200) 'element_id'"`
	CurrencyId  string    `xorm:"varchar(200) 'currency_id'"`
	SrcTypeEnum string    `xorm:"varchar(200) 'src_type_enum'"`
	Money       float64   `xorm:"decimal 'money'"`

	Memo       string `xorm:"varchar(200) 'memo'"`
	StateEnum  string `xorm:"varchar(200) 'state_enum'"`
	ModelingId string `xorm:"varchar(200) 'modeling_id'"`

	CreatedAt time.Time `xorm:"created 'created_at'"`
}

func (s *DataDoc) TableName() string {
	return "suite_amiba_data_docs"
}
