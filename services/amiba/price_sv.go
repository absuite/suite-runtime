package amibaServices

import (
	"errors"
	"fmt"
	"time"

	"github.com/ggoop/goutils/glog"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/absuite/suite-runtime/repositories"
)

type PricelSv interface {
	Cache(entId string, purposeId string) error
	GetPrice(find PriceFind) (amibaModels.Price, error)
}
type priceSv struct {
	repo *repositories.ModelRepo
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

var price_sv *priceSv
var price_sv_cache map[string]map[string]map[string][]amibaModels.Price

func NewPricelSv(repo *repositories.ModelRepo) PricelSv {
	price_sv = &priceSv{repo: repo}
	price_sv_cache = make(map[string]map[string]map[string][]amibaModels.Price)
	return price_sv
}
func (s *priceSv) getPriceInList(find PriceFind, items []amibaModels.Price) (amibaModels.Price, bool) {
	founds := [3]amibaModels.Price{}
	for _, p := range items {
		if find.PriceId != "" && p.Id != find.PriceId { //按价表查询，非此价表的不查
			continue
		}
		if !find.Date.IsZero() && (find.Date.Unix() < p.FmDate.Unix() || find.Date.Unix() > p.ToDate.Unix()) { //数据在日期范围之外，则直接返回
			continue
		}
		if p.FmGroupId == find.FmGroupId && p.ToGroupId == find.ToGroupId && p.ItemCode == find.ItemCode { //来源巴+去向巴+物料,如果匹配，则直接返回
			founds[0] = p
			break
		}
		if p.FmGroupId == find.FmGroupId && p.ToGroupId == "" && p.ItemCode == find.ItemCode { //来源巴+物料
			founds[1] = p
			continue
		}
		if p.FmGroupId == "" && p.ToGroupId == "" && p.ItemCode == find.ItemCode { //物料
			founds[2] = p
			continue
		}
	}
	if founds[0].Id != "" {
		return founds[0], true
	} else if founds[1].Id != "" {
		return founds[1], true
	} else if founds[2].Id != "" {
		return founds[2], true
	}
	return amibaModels.Price{}, false
}
func (s *priceSv) GetPrice(find PriceFind) (amibaModels.Price, error) {
	price := amibaModels.Price{CostPrice: 0}
	items := s.GetCache(find.EntId, find.PurposeId)
	if find.ItemCode == "" || find.FmGroupId == "" || items == nil {
		return price, errors.New("没有找到价格!")
	}
	//1.来源巴+去向巴+物料
	if find.FmGroupId != "" && find.ToGroupId != "" && find.ItemCode != "" {
		findKey := s.getItemKey(find.FmGroupId, find.ToGroupId, find.ItemCode)
		price, found := s.getPriceInList(find, items[findKey])
		if found {
			return price, nil
		}
	}
	//2.来源巴+物料
	if find.FmGroupId != "" && find.ItemCode != "" {
		findKey := s.getItemKey(find.FmGroupId, "", find.ItemCode)
		price, found := s.getPriceInList(find, items[findKey])
		if found {
			return price, nil
		}
	}
	//3.物料
	if find.ItemCode != "" {
		findKey := s.getItemKey("", "", find.ItemCode)
		price, found := s.getPriceInList(find, items[findKey])
		if found {
			return price, nil
		}
	}
	return price, errors.New("没有找到价格!")
}
func (s *priceSv) GetCache(entId string, purposeId string) map[string][]amibaModels.Price {
	e := price_sv_cache[entId]
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
	return nil
}
func (s *priceSv) Cache(entId string, purposeId string) error {
	datas := make([]amibaModels.Price, 0)

	items := make([]amibaModels.Price, 0)
	query := price_sv.repo.Select(`
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
	query = price_sv.repo.Select(`
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
		if item.FmGroupId != "" && item.ToGroupId != "" && item.ItemCode != "" { //来源巴+去向巴+物料,如果匹配，则直接返回
			cacheKey = s.getItemKey(item.FmGroupId, item.ToGroupId, item.ItemCode)
		}
		if item.FmGroupId != "" && item.ToGroupId == "" && item.ItemCode != "" { //来源巴+物料
			cacheKey = s.getItemKey(item.FmGroupId, item.ToGroupId, item.ItemCode)
		}
		if item.FmGroupId == "" && item.ToGroupId == "" && item.ItemCode != "" { //物料
			cacheKey = s.getItemKey(item.FmGroupId, item.ToGroupId, item.ItemCode)
		}
		if cacheKey == "" {
			continue
		}
		if cacheItems[cacheKey] == nil {
			cacheItems[cacheKey] = make([]amibaModels.Price, 0)
		}
		cacheItems[cacheKey] = append(cacheItems[cacheKey], item)
	}
	if price_sv_cache[entId] == nil {
		price_sv_cache[entId] = make(map[string]map[string][]amibaModels.Price)
	}
	price_sv_cache[entId][purposeId] = cacheItems
	return nil
}
