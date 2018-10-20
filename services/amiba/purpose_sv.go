package amibaServices

import (
	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type PurposeSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (amibaModels.Purpose, bool)
	FindByCode(entId string, code string) (amibaModels.Purpose, bool)
}
type purposeSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]amibaModels.Purpose
	cached map[string]bool
}

func NewPurposeSv(repo *repositories.ModelRepo) PurposeSv {
	return &purposeSv{repo: repo, cache: make(map[string]map[string]amibaModels.Purpose), cached: make(map[string]bool)}

}
func (s *purposeSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *purposeSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]amibaModels.Purpose)
	items := make([]amibaModels.Purpose, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_amiba_purposes").Alias("d")
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
func (s *purposeSv) beforeCacheGet(entId string) {
	if !s.cached[entId] {
		s.Cache(entId)
	}
}
func (s *purposeSv) Get(entId string, id string) (amibaModels.Purpose, bool) {
	if entId == "" || id == "" {
		return amibaModels.Purpose{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *purposeSv) FindByCode(entId string, code string) (amibaModels.Purpose, bool) {
	if entId == "" || code == "" {
		return amibaModels.Purpose{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
