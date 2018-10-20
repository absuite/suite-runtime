package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type EntSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(id string) (cboModels.Ent, bool)
	FindByCode(code string) (cboModels.Ent, bool)
}
type entSv struct {
	repo  *repositories.ModelRepo  `optional:"false"`
	cache map[string]cboModels.Ent `optional:"false"`
}

func NewEntSv(repo *repositories.ModelRepo) EntSv {
	return &entSv{repo: repo, cache: make(map[string]cboModels.Ent)}
}
func (s *entSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *entSv) Cache(entId string) error {
	items := make([]cboModels.Ent, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_cbo_orgs").Alias("d")
	query.Where("d.ent_id = ? ", entId)
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	for _, item := range items {
		s.cache["id:"+item.Id] = item
		s.cache["code:"+item.Id] = item
	}
	return nil
}
func (s *entSv) Get(id string) (cboModels.Ent, bool) {
	v, f := s.cache["id:"+id]
	return v, f
}
func (s *entSv) FindByCode(code string) (cboModels.Ent, bool) {
	v, f := s.cache["code:"+code]
	return v, f
}
