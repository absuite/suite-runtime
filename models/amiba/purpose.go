package amibaModels

type Purpose struct {
	ID         string `xorm:"varchar(200) 'id'"`
	Code       string `xorm:"varchar(200) 'code'"`
	Name       string `xorm:"varchar(200) 'name'"`
	CalendarId string `xorm:"varchar(200) 'calendar_id'"`
}
