package amibaModels

import "time"

type DataBiz struct {
	Id           string    `xorm:"varchar(200) 'id'"`
	EntId        string    `xorm:"varchar(200) 'ent_id'"`
	DocNo        string    `xorm:"varchar(200) 'doc_no'"`
	DocDate      time.Time `xorm:"timestamp 'doc_date'"`
	SrcDocId     string    `xorm:"varchar(200) 'src_doc_id'"`
	SrcDocType   string    `xorm:"varchar(200) 'src_doc_type'"`
	BizType      string    `xorm:"varchar(200) 'biz_type'"`
	DocType      string    `xorm:"varchar(200) 'doc_type'"`
	Org          string    `xorm:"varchar(200) 'org'"`
	Person       string    `xorm:"varchar(200) 'person'"`
	FmOrg        string    `xorm:"varchar(200) 'fm_org'"`
	FmDept       string    `xorm:"varchar(200) 'fm_dept'"`
	FmWork       string    `xorm:"varchar(200) 'fm_work'"`
	FmTeam       string    `xorm:"varchar(200) 'fm_team'"`
	FmWh         string    `xorm:"varchar(200) 'fm_wh'"`
	FmPerson     string    `xorm:"varchar(200) 'fm_person'"`
	ToOrg        string    `xorm:"varchar(200) 'to_org'"`
	ToDept       string    `xorm:"varchar(200) 'to_dept'"`
	ToWork       string    `xorm:"varchar(200) 'to_work'"`
	ToTeam       string    `xorm:"varchar(200) 'to_team'"`
	ToWh         string    `xorm:"varchar(200) 'to_wh'"`
	ToPerson     string    `xorm:"varchar(200) 'to_person'"`
	Direction    string    `xorm:"varchar(200) 'direction'"`
	Trader       string    `xorm:"varchar(200) 'trader'"`
	Item         string    `xorm:"varchar(200) 'item'"`
	ItemCategory string    `xorm:"varchar(200) 'item_category'"`
	Project      string    `xorm:"varchar(200) 'project'"`
	Account      string    `xorm:"varchar(200) 'account'"`
	Currency     string    `xorm:"varchar(200) 'currency'"`
	Uom          string    `xorm:"varchar(200) 'uom'"`
	Qty          float64   `xorm:"decimal 'qty'"`
	Price        float64   `xorm:"decimal 'price'"`
	Money        float64   `xorm:"decimal 'money'"`
	Tax          float64   `xorm:"decimal 'tax'"`
	CreditMoney  float64   `xorm:"decimal 'credit_money'"`
	DebitMoney   float64   `xorm:"decimal 'debit_money'"`

	Factor1         string `xorm:"varchar(200) 'factor1'"`
	Factor2         string `xorm:"varchar(200) 'factor2'"`
	Factor3         string `xorm:"varchar(200) 'factor3'"`
	Factor4         string `xorm:"varchar(200) 'factor4'"`
	Factor5         string `xorm:"varchar(200) 'factor5'"`
	DataSrcIdentity string `xorm:"varchar(200) 'data_src_identity'"`
}
