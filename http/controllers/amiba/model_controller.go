package amibaControllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/absuite/suite-runtime/results"
	"github.com/absuite/suite-runtime/services/amiba"
	"github.com/absuite/suite-runtime/services/cbo"
	"github.com/ggoop/goutils/glog"
	"github.com/kataras/iris"
)

type inputModling struct {
	EntId     string   `json:"ent_id"`
	PurposeId string   `json:"purpose_id"`
	PeriodIds []string `json:"period_ids"`
	ModelIds  []string `json:"model_ids"`
	Memo      string   `json:"memo"`
}
type ModelController struct {
	Ctx       iris.Context
	ModelSv   amibaServices.ModelSv
	PurposeSv amibaServices.PurposeSv
	EntSv     cboServices.EntSv
	PeriodSv  cboServices.PeriodSv
}

func (c *ModelController) GetModelingTest() results.Result {
	m := inputModling{}
	m.PurposeId = c.Ctx.URLParam("purpose_id")
	if v := c.Ctx.URLParam("period_ids"); v != "" {
		m.PeriodIds = []string{v}
	}
	if v := c.Ctx.URLParam("model_ids"); v != "" {
		m.ModelIds = []string{v}
	}
	m.EntId = c.Ctx.URLParam("ent_id")
	return c.handModeling(m)
}
func (c *ModelController) handModeling(m inputModling) results.Result {
	fm_time := time.Now()
	ent, f := c.EntSv.Get(m.EntId)
	if !f {
		return results.ToError(errors.New(fmt.Sprintf("企业参数错误:%s", m.EntId)))
	}
	if m.PeriodIds == nil || len(m.PeriodIds) == 0 {
		return results.ToError(errors.New("缺少期间参数!"))
	}
	purpose, f := c.PurposeSv.Get(ent.Id, m.PurposeId)
	if !f {
		return results.ToError(errors.New(fmt.Sprintf("核算目的参数错误:%s", m.PurposeId)))
	}
	result := make(map[string]interface{})
	if m.PeriodIds != nil && len(m.PeriodIds) > 0 {
		for _, periodId := range m.PeriodIds {
			period, f := c.PeriodSv.Get(ent.Id, periodId)
			if !f {
				err := errors.New(fmt.Sprintf("企业:%v,找不到期间数据:%s", ent.Name, periodId))
				glog.Printf("period data error :%s", err)
				continue
			}
			res, err := c.ModelSv.Modeling(ent, purpose, period, m.ModelIds)
			if err != nil {
				result[periodId] = err.Error()
			} else {
				result[periodId] = res
			}
		}
	}
	glog.Printf("企业:%v,核算目的:%v,期间:%v,模型:%v,time:%v Seconds", ent.Name, purpose.Name, m.PeriodIds, m.ModelIds, time.Now().Sub(fm_time).Seconds())
	return results.ToJson(results.Map{"data": result})
}
func (c *ModelController) PostModeling() results.Result {
	var m inputModling
	if err := c.Ctx.ReadJSON(&m); err != nil {
		return results.ToError(err)
	}
	m.EntId = c.Ctx.Values().GetString("Ent")
	return c.handModeling(m)
}
