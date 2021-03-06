package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type ItemSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Item, bool)
	FindByCode(entId string, code string) (cboModels.Item, bool)
}
type itemSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]cboModels.Item
	cached map[string]bool
}

func NewItemSv(repo *repositories.ModelRepo) ItemSv {
	return &itemSv{repo: repo, cache: make(map[string]map[string]cboModels.Item), cached: make(map[string]bool)}
}

func (s *itemSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *itemSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.Item)
	items := make([]cboModels.Item, 0)
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
func (s *itemSv) beforeCacheGet(entId string) {
	if !s.cached[entId] {
		s.Cache(entId)
	}
}
func (s *itemSv) Get(entId string, id string) (cboModels.Item, bool) {
	if entId == "" || id == "" {
		return cboModels.Item{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *itemSv) FindByCode(entId string, code string) (cboModels.Item, bool) {
	if entId == "" || code == "" {
		return cboModels.Item{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
