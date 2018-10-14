package cboModels

type Org struct {
	Id    string `xorm:"'id'"`
	EntId string `xorm:"varchar(200) 'ent_id'"`
	Code  string `xorm:"'code'"`
	Name  string `xorm:"'name'"`
}
