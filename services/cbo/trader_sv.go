package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type TraderSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Trader, bool)
	FindByCode(entId string, code string) (cboModels.Trader, bool)
}
type traderSvSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]cboModels.Trader
	cached map[string]bool
}

func NewTraderSv(repo *repositories.ModelRepo) TraderSv {
	return &traderSvSv{repo: repo, cache: make(map[string]map[string]cboModels.Trader), cached: make(map[string]bool)}

}
func (s *traderSvSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *traderSvSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.Trader)
	items := make([]cboModels.Trader, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_cbo_traders").Alias("d")
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
func (s *traderSvSv) beforeCacheGet(entId string) {
	if !s.cached[entId] {
		s.Cache(entId)
	}
}
func (s *traderSvSv) Get(entId string, id string) (cboModels.Trader, bool) {
	if entId == "" || id == "" {
		return cboModels.Trader{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *traderSvSv) FindByCode(entId string, code string) (cboModels.Trader, bool) {
	if entId == "" || code == "" {
		return cboModels.Trader{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
