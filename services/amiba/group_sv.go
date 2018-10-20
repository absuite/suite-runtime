package amibaServices

import (
	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/ggoop/goutils/glog"
)

type GroupSv interface {
	Cache(entId string) error
	CacheAll() error
	Get(entId string, id string) (amibaModels.Group, bool)
	FindByCode(entId, purposeId, code string) (amibaModels.Group, bool)
	FindByLineCode(entId, purposeId, code string) (amibaModels.Group, bool)
}
type groupSv struct {
	repo   *repositories.ModelRepo
	cache  map[string]map[string]amibaModels.Group
	cached map[string]bool
}

func NewGroupSv(repo *repositories.ModelRepo) GroupSv {
	return &groupSv{repo: repo, cache: make(map[string]map[string]amibaModels.Group), cached: make(map[string]bool)}

}
func (s *groupSv) CacheAll() error {
	ents := s.repo.GetEnts()
	for _, v := range ents {
		s.Cache(v.Id)
	}
	return nil
}
func (s *groupSv) Cache(entId string) error {
	s.cache[entId] = make(map[string]amibaModels.Group)
	items := make([]amibaModels.Group, 0)
	query := s.repo.Select("g.id,g.code,g.name,g.type_enum").Table("suite_amiba_groups").Alias("g")
	query.Where("g.ent_id = ? ", entId)
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	groupData := make([]amibaModels.GroupData, 0)
	groupDataItem := make([]amibaModels.GroupData, 0)

	query = s.repo.Select("g.id as group_id,g.type_enum as type_enum,d.id,d.code,d.name").Table("suite_amiba_groups").Alias("g")
	query.Join("inner", []string{"suite_amiba_group_lines", "gl"}, "g.id=gl.group_id")
	query.Join("inner", []string{"suite_cbo_orgs", "d"}, "gl.data_id=d.id").Where("g.type_enum=?", "org")
	if err := query.Find(&groupDataItem); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if len(groupDataItem) > 0 {
		groupData = append(groupData, groupDataItem...)
	}

	query = s.repo.Select("g.id as group_id,g.type_enum as type_enum,d.id,d.code,d.name").Table("suite_amiba_groups").Alias("g")
	query.Join("inner", []string{"suite_amiba_group_lines", "gl"}, "g.id=gl.group_id")
	query.Join("inner", []string{"suite_cbo_depts", "d"}, "gl.data_id=d.id").Where("g.type_enum=?", "dept")
	if err := query.Find(&groupDataItem); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if len(groupDataItem) > 0 {
		groupData = append(groupData, groupDataItem...)
	}

	query = s.repo.Select("g.id as group_id,g.type_enum as type_enum,d.id,d.code,d.name").Table("suite_amiba_groups").Alias("g")
	query.Join("inner", []string{"suite_amiba_group_lines", "gl"}, "g.id=gl.group_id")
	query.Join("inner", []string{"suite_cbo_works", "d"}, "gl.data_id=d.id").Where("g.type_enum=?", "work")
	if err := query.Find(&groupDataItem); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if len(groupDataItem) > 0 {
		groupData = append(groupData, groupDataItem...)
	}
	if len(groupData) > 0 {
		for gi, gv := range items {
			for _, vv := range groupData {
				if gv.Id == vv.GroupId {
					items[gi].AddData(vv)
				}
			}
		}
	}
	for _, item := range items {
		s.cache[entId]["id:"+item.Id] = item
		s.cache[entId]["purpose:"+item.PurposeId+"code:"+item.Code] = item
		if len(item.Datas) > 0 {
			for _, line := range item.Datas {
				s.cache[entId]["purpose:"+item.PurposeId+"code:line:"+line.Code] = item
			}
		}
	}
	s.cached[entId] = true
	return nil
}
func (s *groupSv) beforeCacheGet(entId string) {
	if !s.cached[entId] {
		s.Cache(entId)
	}
}
func (s *groupSv) Get(entId string, id string) (amibaModels.Group, bool) {
	if entId == "" || id == "" {
		return amibaModels.Group{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["id:"+id]
	return v, f
}
func (s *groupSv) FindByCode(entId, purposeId, code string) (amibaModels.Group, bool) {
	if entId == "" || purposeId == "" || code == "" {
		return amibaModels.Group{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["purpose:"+purposeId+"code:"+code]
	return v, f
}
func (s *groupSv) FindByLineCode(entId, purposeId, code string) (amibaModels.Group, bool) {
	if entId == "" || purposeId == "" || code == "" {
		return amibaModels.Group{}, false
	}
	s.beforeCacheGet(entId)
	v, f := s.cache[entId]["purpose:"+purposeId+"code:line:"+code]
	return v, f
}
