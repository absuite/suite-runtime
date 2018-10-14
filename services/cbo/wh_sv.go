package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type WhSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Wh, bool)
	FindByCode(entId string, code string) (cboModels.Wh, bool)
}
type whSv struct {
	repo  *repositories.ModelRepo
	cache map[string]map[string]cboModels.Wh
}

func NewWhSv(repo *repositories.ModelRepo) WhSv {
	return &whSv{repo: repo, cache: make(map[string]map[string]cboModels.Wh)}

}
func (s *whSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *whSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.Wh)
	items := make([]cboModels.Wh, 0)
	query := s.repo.Select("d.ent_id,d.id,d.code,d.name").Table("suite_cbo_whs").Alias("d")
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
func (s *whSv) Get(entId string, id string) (cboModels.Wh, bool) {
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *whSv) FindByCode(entId string, code string) (cboModels.Wh, bool) {
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
