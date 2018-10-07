package results

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type (
	Result = mvc.Result
	Map    = iris.Map
)

func ToJson(data interface{}) Result {
	return mvc.Response{
		Object: data,
	}
}
func ToError(err error, code ...int) Result {
	if code != nil && len(code) > 0 {
		return mvc.Response{
			Code:   code[0],
			Object: iris.Map{"msg": err.Error()},
		}
	} else {
		return mvc.Response{
			Code:   iris.StatusBadRequest,
			Object: iris.Map{"msg": err.Error()},
		}
	}
}
