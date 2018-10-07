package amibaControllers

import (
	"github.com/absuite/suite-runtime/results"
	"github.com/absuite/suite-runtime/services/amiba"
	"github.com/kataras/iris"
)

type priceInput struct {
	EntId     string `json:"ent_id"`
	PurposeId string `json:"purpose_id"`
}
type PriceController struct {
	Ctx     iris.Context
	Service amibaServices.PricelSv
}

func (c *PriceController) PostCache() results.Result {
	var m priceInput
	if err := c.Ctx.ReadJSON(&m); err != nil {
		return results.ToError(err)
	}
	m.EntId = c.Ctx.Values().GetString("Ent")
	err := c.Service.Cache(m.EntId, m.PurposeId)
	if err != nil {
		return results.ToError(err)
	} else {
		return results.ToJson(results.Map{"data": true})
	}
}
