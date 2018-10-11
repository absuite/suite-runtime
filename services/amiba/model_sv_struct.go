package amibaServices

import (
	"time"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/models/cbo"
)

type tmlModelingLine struct {
	EntId      string
	Period     cboModels.Period
	ModelLine  amibaModels.ModelLine
	Group      amibaModels.Group
	MatchGroup amibaModels.Group
	AllGroups  []amibaModels.Group
}
type tmlDataElementing struct {
	EntId     string
	PurposeId string
	PeriodId  string
	FmGroupId string //结果来源巴
	ToGroupId string //结果目标巴

	ModelingId         string //模型行ID
	ModelingLineId     string //模型行ID
	MatchDirectionEnum string
	MatchGroupId       string //匹配方：与原始业务数据中的巴进行匹配，并将匹配结果作为表头巴的建模成果。
	DefFmGroupId       string //模型头定义的巴
	DefToGroupId       string //模型头定义的巴,交易方：标识表头巴的交易对方巴是谁，当交易方为空时，使用原始业务数据匹配的巴，如果指定了，则直接使用交易巴。

	ElementId     string
	BizTypeEnum   string //业务类型
	ValueTypeEnum string
	Adjust        string

	DataId        string
	DataType      string
	DataFmGroupId string //业务数据来源对应的巴
	DataToGroupId string //业务数据目标对应的巴

	DataDocNo    string
	DataDocDate  time.Time
	DataFmOrg    string
	DataFmDept   string
	DataFmWork   string
	DataFmTeam   string
	DataFmWh     string
	DataFmPerson string

	DataToOrg    string
	DataToDept   string
	DataToWork   string
	DataToTeam   string
	DataToWh     string
	DataToPerson string

	DataTraderId         string
	DataTraderCode       string
	DataItemCode         string
	DataItemId           string
	DataItemCategoryCode string
	DataItemCategoryId   string
	DataProjectId        string
	DataProjectCode      string
	DataAccountCode      string
	DataCurrency         string
	DataUom              string
	DataQty              float64
	DataMoney            float64

	Qty   float64
	Price float64
	Money float64

	Deleted bool
}
