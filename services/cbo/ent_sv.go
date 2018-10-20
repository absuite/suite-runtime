package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type EntSv interface {
	CacheAll() error
	Get(id string) (cboModels.Ent, bool)
	FindByCode(code string) (cboModels.Ent, bool)
}
type entSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]cboModels.Ent
	cached bool
}

func NewEntSv(repo *repositories.ModelRepo) EntSv {
	return &entSv{repo: repo, cache: make(map[string]cboModels.Ent)}
}
func (s *entSv) CacheAll() error {
	items := make([]cboModels.Ent, 0)
	query := s.repo.Select("d.id,d.code,d.name").Table("gmf_sys_ents").Alias("d")
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	for _, item := range items {
		s.cache["id:"+item.Id] = item
		s.cache["code:"+item.Id] = item
	}
	s.cached = true
	return nil
}
func (s *entSv) beforeCacheGet() {
	if !s.cached {
		s.CacheAll()
	}
}
func (s *entSv) Get(id string) (cboModels.Ent, bool) {
	if id == "" {
		return cboModels.Ent{}, false
	}
	s.beforeCacheGet()
	v, f := s.cache["id:"+id]
	return v, f
}
func (s *entSv) FindByCode(code string) (cboModels.Ent, bool) {
	if code == "" {
		return cboModels.Ent{}, false
	}
	s.beforeCacheGet()
	v, f := s.cache["code:"+code]
	return v, f
}
