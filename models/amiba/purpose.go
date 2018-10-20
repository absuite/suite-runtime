package amibaModels

type Purpose struct {
	Id    string `xorm:"varchar(200) 'id'"`
	EntId string `xorm:"varchar(200) 'ent_id'"`

	Code       string `xorm:"varchar(200) 'code'"`
	Name       string `xorm:"varchar(200) 'name'"`
	CalendarId string `xorm:"varchar(200) 'calendar_id'"`
}
