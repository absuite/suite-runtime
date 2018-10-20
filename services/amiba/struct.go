package amibaServices

import (
	"time"
)

type IDCodeName struct {
	Id   string `xorm:"varchar(200) 'id'"`
	Code string `xorm:"varchar(200) 'code'"`
	Name string `xorm:"varchar(200) 'name'"`
}

var fm_time time.Time
