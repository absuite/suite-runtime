package repositories

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/ggoop/goutils/glog"
	"github.com/go-xorm/xorm"
)

type ModelRepo struct {
	*xorm.Engine
}

var model_repo_cache []cboModels.Ent

func NewModelRepo(orm *xorm.Engine) *ModelRepo {
	return &ModelRepo{orm}
}
func (s *ModelRepo) GetEnts() []cboModels.Ent {
	if model_repo_cache != nil && len(model_repo_cache) > 0 {
		return model_repo_cache
	}
	model_repo_cache = make([]cboModels.Ent, 0)
	query := s.Select(`p.id,p.code,p.name,p.ent_id,p.calendar_id`)
	query.Table("gmf_sys_ents").Alias("p")
	if err := query.Find(&model_repo_cache); err != nil {
		glog.Printf("query ents error :%s", err)
	}
	return model_repo_cache
}
func (s *ModelRepo) GetEnt(entId string) (cboModels.Ent, bool) {
	if model_repo_cache != nil && len(model_repo_cache) >= 0 {
		for _, v := range model_repo_cache {
			if entId == v.Id {
				return v, true
			}
		}
	}
	return cboModels.Ent{}, false
}
