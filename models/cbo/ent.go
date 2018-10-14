package cboModels

type Ent struct {
	Id   string `xorm:"'id'"`
	Code string `xorm:"'code'"`
	Name string `xorm:"'name'"`
}
