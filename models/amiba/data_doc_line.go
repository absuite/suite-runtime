package amibaModels

import "time"

type DataDocLine struct {
	Id    string `xorm:"varchar(200) 'id'"`
	EntId string `xorm:"varchar(200) 'ent_id'"`
	DocId string `xorm:"varchar(200) 'doc_id'"`

	TraderId       string  `xorm:"varchar(200) 'trader_id'"`
	ItemCategoryId string  `xorm:"varchar(200) 'item_category_id'"`
	ItemId         string  `xorm:"varchar(200) 'item_id'"`
	MfcId          string  `xorm:"varchar(200) 'mfc_id'"`
	ProjectId      string  `xorm:"varchar(200) 'project_id'"`
	ExpenseCode    string  `xorm:"varchar(200) 'expense_code'"`
	AccountCode    string  `xorm:"varchar(200) 'account_code'"`
	UnitId         string  `xorm:"varchar(200) 'unit_id'"`
	Qty            float64 `xorm:"decimal 'qty'"`
	Price          float64 `xorm:"decimal 'price'"`
	Money          float64 `xorm:"decimal 'money'"`

	ModelingId     string    `xorm:"varchar(200) 'modeling_id'"`
	ModelingLineId string    `xorm:"varchar(200) 'modeling_line_id'"`
	CreatedAt      time.Time `xorm:"created 'created_at'"`
}
