package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type ProjectSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Project, bool)
	FindByCode(entId string, code string) (cboModels.Project, bool)
}
type projectSv struct {
	repo  *repositories.ModelRepo
	cache map[string]map[string]cboModels.Project
}

func NewProjectSv(repo *repositories.ModelRepo) ProjectSv {
	return &projectSv{repo: repo, cache: make(map[string]map[string]cboModels.Project)}

}
func (s *projectSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *projectSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.Project)
	items := make([]cboModels.Project, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_cbo_projects").Alias("d")
	query.Where("d.ent_id = ? ", entId)
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	for _, item := range items {
		s.cache[entId]["id:"+item.Id] = item
		s.cache[entId]["code:"+item.Id] = item
	}
	return nil
}
func (s *projectSv) Get(entId string, id string) (cboModels.Project, bool) {
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *projectSv) FindByCode(entId string, code string) (cboModels.Project, bool) {
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
