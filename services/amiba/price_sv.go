package amibaServices

import (
	"fmt"
	"time"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/repositories"
)

type PricelSv interface {
	Cache(entId string, purposeId string) error
	GetPrice(find PriceFind) (amibaModels.Price, bool, error)
	CacheAll() error
}
type priceSv struct {
	repo  *repositories.ModelRepo
	cache map[string]map[string]map[string][]amibaModels.Price
}
type PriceFind struct {
	EntId     string
	PurposeId string
	PriceId   string
	FmGroupId string
	ToGroupId string
	ItemCode  string
	Date      time.Time
}

func NewPricelSv(repo *repositories.ModelRepo) PricelSv {
	return &priceSv{repo: repo, cache: make(map[string]map[string]map[string][]amibaModels.Price)}
}
func (s *priceSv) getPriceInList(find PriceFind, findKey string, items []amibaModels.Price) (amibaModels.Price, bool) {
	for _, p := range items {
		if find.PriceId != "" && p.Id != find.PriceId { //按价表查询，非此价表的不查
			continue
		}
		if p.CacheKey != findKey {
			continue
		}
		if !find.Date.IsZero() && (find.Date.Unix() < p.FmDate.Unix() || find.Date.Unix() > p.ToDate.Unix()) { //数据在日期范围之外，则直接返回
			continue
		}
		return p, true
	}
	return amibaModels.Price{}, false
}
func (s *priceSv) GetPrice(find PriceFind) (amibaModels.Price, bool, error) {
	price := amibaModels.Price{CostPrice: 0}
	items := s.GetCache(find.EntId, find.PurposeId)
	if items == nil || len(items) == 0 {
		glog.Errorf("企业:%v,核算:%v,找不到价表数据", find.EntId, find.PurposeId)
		return price, false, nil
	}
	if find.ItemCode == "" || find.FmGroupId == "" {
		glog.Errorf("企业:%v,核算:%v,物料:%v,来源巴:%v,参数不符合,不能取价!", find.EntId, find.PurposeId, find.ItemCode, find.FmGroupId)
		return price, false, nil
	}
	//1.来源巴+去向巴+物料
	if find.FmGroupId != "" && find.ToGroupId != "" && find.ItemCode != "" {
		findKey := s.getItemKey(find.FmGroupId, find.ToGroupId, find.ItemCode)
		price, found := s.getPriceInList(find, findKey, items[findKey])
		if found {
			return price, true, nil
		}
	}
	//2.来源巴+物料
	if find.FmGroupId != "" && find.ItemCode != "" {
		findKey := s.getItemKey(find.FmGroupId, "", find.ItemCode)
		price, found := s.getPriceInList(find, findKey, items[findKey])
		if found {
			return price, true, nil
		}
	}
	//3.物料
	if find.ItemCode != "" {
		findKey := s.getItemKey("", "", find.ItemCode)
		price, found := s.getPriceInList(find, findKey, items[findKey])
		if found {
			return price, true, nil
		}
	}
	return price, false, nil
}
func (s *priceSv) GetCache(entId string, purposeId string) map[string][]amibaModels.Price {
	e := s.cache[entId]
	if e == nil || len(e) == 0 {
		return nil
	}
	p := e[purposeId]
	if p == nil || len(p) == 0 {
		return nil
	}
	return p
}
func (s *priceSv) getItemKey(FmGroupId string, ToGroupId string, ItemCode string) string {
	return fmt.Sprintf("%s:%s:%s", FmGroupId, ToGroupId, ItemCode)
}
func (s *priceSv) CacheAll() error {
	items := make([]amibaModels.Purpose, 0)
	query := s.repo.Select(`p.id,p.code,p.name,p.ent_id,p.calendar_id`)
	query.Table("suite_amiba_purposes").Alias("p")
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if items != nil && len(items) > 0 {
		for _, v := range items {
			s.Cache(v.EntId, v.Id)
		}
	}
	return nil
}
func (s *priceSv) Cache(entId string, purposeId string) error {
	datas := make([]amibaModels.Price, 0)

	items := make([]amibaModels.Price, 0)
	query := s.repo.Select(`
		p.id,p.code,p.name,p.purpose_id,p.group_id as fm_group_id,
		pl.fm_date,pl.to_date,pl.group_id as to_group_id,pl.item_id,item.code as item_code,pl.cost_price
	`)
	query.Table("suite_amiba_prices").Alias("p")
	query.Join("inner", []string{"suite_amiba_price_lines", "pl"}, "p.id=pl.price_id")
	query.Join("left", []string{"suite_cbo_items", "item"}, "pl.item_id=item.id")
	query.Where("p.ent_id=? and p.purpose_id=?", entId, purposeId)
	query.Where("(p.disabled=0 or p.disabled is null)")
	query.Desc("p.group_id").Desc("pl.item_id").Desc("pl.group_id").Desc("pl.fm_date")
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if items != nil && len(items) > 0 {
		datas = append(datas, items...)
	}

	items = make([]amibaModels.Price, 0)
	query = s.repo.Select(`
		ph.id,ph.code,ph.name,p.purpose_id,ph.group_id as fm_group_id,
		pl.fm_date,pl.to_date,pl.group_id as to_group_id,pl.item_id,item.code as item_code,pl.cost_price
	`)
	query.Table("suite_amiba_price_adjusts").Alias("p")
	query.Join("inner", []string{"suite_amiba_price_adjust_lines", "pl"}, "p.id=pl.adjust_id")
	query.Join("inner", []string{"suite_amiba_prices", "ph"}, "ph.id=pl.price_id")
	query.Join("left", []string{"suite_cbo_items", "item"}, "pl.item_id=item.id")
	query.Where("p.ent_id=? and p.purpose_id=?", entId, purposeId)
	query.Where("(p.disabled=0 or p.disabled is null) and (ph.disabled=0 or ph.disabled is null)")
	query.Desc("ph.group_id").Desc("pl.item_id").Desc("pl.group_id").Desc("pl.fm_date")
	if err := query.Find(&items); err != nil {
		glog.Printf("query error :%s", err)
		return err
	}
	if items != nil && len(items) > 0 {
		datas = append(datas, items...)
	}
	cacheItems := make(map[string][]amibaModels.Price)
	cacheKey := ""
	for _, item := range datas {
		cacheKey = ""
		if cacheKey == "" && item.FmGroupId != "" && item.ToGroupId != "" && item.ItemCode != "" { //来源巴+去向巴+物料,如果匹配，则直接返回
			cacheKey = s.getItemKey(item.FmGroupId, item.ToGroupId, item.ItemCode)
		}
		if cacheKey == "" && item.FmGroupId != "" && item.ToGroupId == "" && item.ItemCode != "" { //来源巴+物料
			cacheKey = s.getItemKey(item.FmGroupId, item.ToGroupId, item.ItemCode)
		}
		if cacheKey == "" && item.FmGroupId == "" && item.ToGroupId == "" && item.ItemCode != "" { //物料
			cacheKey = s.getItemKey(item.FmGroupId, item.ToGroupId, item.ItemCode)
		}
		if cacheKey == "" {
			continue
		}
		item.CacheKey = cacheKey

		if cacheItems[cacheKey] == nil {
			cacheItems[cacheKey] = make([]amibaModels.Price, 0)
		}
		cacheItems[cacheKey] = append(cacheItems[cacheKey], item)
	}
	if s.cache[entId] == nil {
		s.cache[entId] = make(map[string]map[string][]amibaModels.Price)
	}
	s.cache[entId][purposeId] = cacheItems

	glog.Printf("缓存价表成功 %d 条：ent=%s,purpose=%s", len(datas), entId, purposeId)
	return nil
}
