package repositories

import (
	"github.com/go-xorm/xorm"
)

type ModelRepo struct {
	*xorm.Engine
}

func NewModelRepo(orm *xorm.Engine) *ModelRepo {
	return &ModelRepo{orm}
}
