package cboServices

import (
	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type PeriodSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Period, bool)
}
type periodSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]cboModels.Period
	cached map[string]bool
}

func NewPeriodSv(repo *repositories.ModelRepo) PeriodSv {
	return &periodSv{repo: repo, cache: make(map[string]map[string]cboModels.Period), cached: make(map[string]bool)}

}
func (s *periodSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *periodSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]cboModels.Period)
	items := make([]cboModels.Period, 0)
	query := s.repo.Select("c.type_enum,c.id as calendar_id,c.name as calendar_name,p.id,p.code,p.name,p.year,p.from_date,p.to_date").Table("suite_cbo_period_calendars").Alias("c")
	query.Join("inner", []string{"suite_cbo_period_accounts", "p"}, "c.id=p.calendar_id")
	query.Where("p.ent_id = ? ", entId)
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	for _, item := range items {
		s.cache[entId]["id:"+item.Id] = item
	}
	s.cached[entId] = true
	return nil
}
func (s *periodSv) beforeCacheGet(entId string) {
	if !s.cached[entId] {
		s.Cache(entId)
	}
}
func (s *periodSv) Get(entId string, id string) (cboModels.Period, bool) {
	if entId == "" || id == "" {
		return cboModels.Period{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
