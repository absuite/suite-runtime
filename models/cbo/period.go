package cboModels

import "time"

type Period struct {
	CalendarId string    `xorm:"varchar(200)  'calendar_id'"`
	Type       string    `xorm:"varchar(200)  'type_enum'"`
	Id         string    `xorm:"'id'"`
	Code       string    `xorm:"'code'"`
	Name       string    `xorm:"'name'"`
	Year       uint      `xorm:"'year'"`
	FromDate   time.Time `xorm:"varchar(200)  'from_date'"`
	ToDate     time.Time `xorm:"varchar(200)  'to_date'"`
}

func (s *Period) TableName() string {
	return "suite_cbo_period_accounts"
}
