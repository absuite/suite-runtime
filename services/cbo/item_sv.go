package cboServices

import (
	"time"

	"github.com/absuite/suite-runtime/models/cbo"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type ItemSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (cboModels.Item, bool)
	FindByCode(entId string, code string) (cboModels.Item, bool)
	GetTestId() int64
}
type itemSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]cboModels.Item
	testId int64
}

func NewItemSv(repo *repositories.ModelRepo) ItemSv {
	return &itemSv{repo: repo, cache: make(map[string]map[string]cboModels.Item), testId: time.Now().UnixNano()}
}
func (s *itemSv) GetTestId() int64 {
	return s.testId
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
	return nil
}
func (s *itemSv) Get(entId string, id string) (cboModels.Item, bool) {
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *itemSv) FindByCode(entId string, code string) (cboModels.Item, bool) {
	v, f := s.cache[entId]["code:"+code]
	return v, f
}
