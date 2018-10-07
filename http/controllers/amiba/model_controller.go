package amibaControllers

import (
	"errors"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/results"
	"github.com/absuite/suite-runtime/services/amiba"
	"github.com/kataras/iris"
)

type modelModling struct {
	PurposeId string   `json:"purpose_id"`
	PeriodIds []string `json:"period_ids"`
	ModelIds  []string `json:"model_ids"`
	Memo      string   `json:"memo"`
}
type ModelController struct {
	Ctx iris.Context
	// Our MovieService, it's an interface which
	// is binded from the main application.
	Service amibaServices.ModelSv
}

func (c *ModelController) PostModeling() results.Result {
	var m modelModling
	if err := c.Ctx.ReadJSON(&m); err != nil {
		return results.ToError(err)
	}
	entId := c.Ctx.Values().GetString("Ent")
	if m.PeriodIds == nil || len(m.PeriodIds) == 0 {
		return results.ToError(errors.New("缺少期间参数!"))
	}
	result := make(map[string]interface{})
	if m.ModelIds != nil && len(m.ModelIds) > 0 {
		for _, modelId := range m.ModelIds {
			for _, periodId := range m.PeriodIds {
				res, err := c.Service.Modeling(amibaModels.Modeling{EntId: entId, PurposeId: m.PurposeId, PeriodId: periodId, ModelId: modelId})
				if err != nil {
					result[modelId+":"+periodId] = err.Error()
				} else {
					result[modelId+":"+periodId] = res
				}
			}
		}
	} else {
		for _, periodId := range m.PeriodIds {
			res, err := c.Service.Modeling(amibaModels.Modeling{EntId: entId, PurposeId: m.PurposeId, PeriodId: periodId})
			if err != nil {
				result[periodId] = err.Error()
			} else {
				result[periodId] = res
			}
		}
	}
	return results.ToJson(results.Map{"data": result})
}
