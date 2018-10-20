package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type DocTypeSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.DocType, bool)
	FindByCode(entId string, code string) (cboModels.DocType, bool)
}
type docTypeSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]cboModels.DocType
	cached map[string]bool
}

func NewDocTypeSv(repo *repositories.ModelRepo) DocTypeSv {
	return &docTypeSv{repo: repo, cache: make(map[string]map[string]cboModels.DocType), cached: make(map[string]bool)}

}
func (s *docTypeSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *docTypeSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.DocType)
	items := make([]cboModels.DocType, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_cbo_doc_types").Alias("d")
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
func (s *docTypeSv) beforeCacheGet(entId string) {
	if !s.cached[entId] {
		s.Cache(entId)
	}
}
func (s *docTypeSv) Get(entId string, id string) (cboModels.DocType, bool) {
	if entId == "" || id == "" {
		return cboModels.DocType{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *docTypeSv) FindByCode(entId string, code string) (cboModels.DocType, bool) {
	if entId == "" || code == "" {
		return cboModels.DocType{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
