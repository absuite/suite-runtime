package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type OrgSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Org, bool)
	FindByCode(entId string, code string) (cboModels.Org, bool)
}
type orgSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]cboModels.Org
	cached map[string]bool
}

func NewOrgSv(repo *repositories.ModelRepo) OrgSv {
	return &orgSv{repo: repo, cache: make(map[string]map[string]cboModels.Org), cached: make(map[string]bool)}
}
func (s *orgSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *orgSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.Org)

	items := make([]cboModels.Org, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_cbo_orgs").Alias("d")
	query.Where("d.ent_id = ? ", entId)
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	for _, item := range items {
		s.cache[entId]["id:"+item.Id] = item
		s.cache[entId]["code:"+item.Id] = item
	}
	s.cached[entId] = true
	return nil
}
func (s *orgSv) beforeCacheGet(entId string) {
	if !s.cached[entId] {
		s.Cache(entId)
	}
}
func (s *orgSv) Get(entId string, id string) (cboModels.Org, bool) {
	if entId == "" || id == "" {
		return cboModels.Org{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *orgSv) FindByCode(entId string, code string) (cboModels.Org, bool) {
	if entId == "" || code == "" {
		return cboModels.Org{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
