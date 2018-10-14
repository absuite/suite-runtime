package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type DeptSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Dept, bool)
	FindByCode(entId string, code string) (cboModels.Dept, bool)
}
type deptSv struct {
	repo  *repositories.ModelRepo
	cache map[string]map[string]cboModels.Dept
}

func NewDeptSv(repo *repositories.ModelRepo) DeptSv {
	return &deptSv{repo: repo, cache: make(map[string]map[string]cboModels.Dept)}

}
func (s *deptSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *deptSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.Dept)
	items := make([]cboModels.Dept, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_cbo_depts").Alias("d")
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
func (s *deptSv) Get(entId string, id string) (cboModels.Dept, bool) {
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *deptSv) FindByCode(entId string, code string) (cboModels.Dept, bool) {
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
