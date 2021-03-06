package amibaServices

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/absuite/suite-runtime/models/cbo"

	"github.com/absuite/suite-runtime/models/amiba"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func (s *modelSv) Savedoc(ent cboModels.Ent, purpose amibaModels.Purpose, period cboModels.Period, items []tmlDataElementing, m amibaModels.Model) {
	fm_time = time.Now()
	var sql string
	if m.Id != "" {
		sql = `
		DELETE l FROM suite_amiba_data_doc_lines AS l	INNER JOIN suite_amiba_data_docs AS h ON l.doc_id=h.id
		WHERE h.src_type_enum=? AND h.ent_id=? AND h.purpose_id=? AND h.period_id=? AND h.modeling_id=?
		`
		if _, err := s.repo.Exec(sql, "interface", ent.Id, m.PurposeId, period.Id, m.Id); err != nil {
			glog.CheckAndPrintError("delete suite_amiba_data_doc_lines error!", err)
		}
	} else {
		sql = `
		DELETE l FROM suite_amiba_data_doc_lines AS l	INNER JOIN suite_amiba_data_docs AS h ON l.doc_id=h.id
		WHERE h.src_type_enum=? AND h.ent_id=? AND h.purpose_id=? AND h.period_id=? 
		`
		if _, err := s.repo.Exec(sql, "interface", ent.Id, m.PurposeId, period.Id); err != nil {
			glog.CheckAndPrintError("delete suite_amiba_data_doc_lines error!", err)
		}
	}
	if m.Id != "" {
		sql = `
		DELETE h FROM suite_amiba_data_docs AS h
		WHERE h.src_type_enum=? AND h.ent_id=? AND h.purpose_id=? AND h.period_id=? AND h.modeling_id=?
		`
		if _, err := s.repo.Exec(sql, "interface", ent.Id, m.PurposeId, period.Id, m.Id); err != nil {
			glog.CheckAndPrintError("delete suite_amiba_data_docs error!", err)
		}
	} else {
		sql = `
		DELETE h FROM suite_amiba_data_docs AS h
		WHERE h.src_type_enum=? AND h.ent_id=? AND h.purpose_id=? AND h.period_id=?
		`
		if _, err := s.repo.Exec(sql, "interface", ent.Id, m.PurposeId, period.Id); err != nil {
			glog.CheckAndPrintError("delete suite_amiba_data_docs error!", err)
		}
	}
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,删除上次建模数据,time:%v Seconds", ent.Name, purpose.Name, m.Name, period.Name, time.Now().Sub(fm_time).Seconds())

	if items == nil || len(items) == 0 {
		return
	}

	fm_time = time.Now()
	groups := make(map[string][]tmlDataElementing)
	groupKey := ""
	for _, item := range items {
		groupKey = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%s:%s:%s", item.ModelingId, item.PurposeId, item.FmGroupId, item.ToGroupId, item.ElementId))))
		if groups[groupKey] == nil {
			groups[groupKey] = make([]tmlDataElementing, 0)
		}
		groups[groupKey] = append(groups[groupKey], item)
	}

	dataDocs := make([]amibaModels.DataDoc, 0)
	dataDocLines := make([]amibaModels.DataDocLine, 0)

	var doc amibaModels.DataDoc
	var docLine amibaModels.DataDocLine
	var count = 0
	var totalMoney float64
	for _, gv := range groups {
		count++
		totalMoney = 0
		for ik, iv := range gv {
			totalMoney += iv.Money
			if ik == 0 {
				doc = amibaModels.DataDoc{EntId: iv.EntId, ModelingId: iv.ModelingId, PurposeId: iv.PurposeId, PeriodId: iv.PeriodId, SrcTypeEnum: "interface"}
				doc.Id = utils.GUID()
				doc.DocNo = time.Now().Format("060102") + utils.NewStringInt(count).ToString()
				doc.DocDate = time.Now()
				doc.FmGroupId = iv.FmGroupId
				doc.ToGroupId = iv.ToGroupId
				doc.ElementId = iv.ElementId
				doc.StateEnum = "approved"
			}
			docLine = amibaModels.DataDocLine{DocId: doc.Id, EntId: iv.EntId, ModelingId: iv.ModelingId, ModelingLineId: iv.ModelingLineId}
			docLine.Id = utils.GUID()
			docLine.Qty = iv.Qty
			if iv.Qty != 0 {
				docLine.Price = iv.Money / iv.Qty
			}
			docLine.Money = iv.Money
			docLine.TraderId = iv.DataTraderId
			docLine.ItemCategoryId = iv.DataItemCategoryId
			docLine.ItemId = iv.DataItemId
			docLine.AccountCode = iv.DataAccountCode
			docLine.ProjectId = iv.DataProjectId

			dataDocLines = append(dataDocLines, docLine)
		}
		doc.Money = totalMoney
		dataDocs = append(dataDocs, doc)
	}
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,处理数据分组完成:原始数据%v条,time:%v Seconds", ent.Name, purpose.Name, m.Name, period.Name, len(groups), time.Now().Sub(fm_time).Seconds())

	//分批插入头
	fm_time = time.Now()
	batchDocs := make([]amibaModels.DataDoc, 0)
	count = 0
	for _, item := range dataDocs {
		count++
		batchDocs = append(batchDocs, item)
		if count > 1000 {
			_, err := s.repo.Table("suite_amiba_data_docs").Insert(&batchDocs)
			if err != nil {
				glog.CheckAndPrintError("insert into docs error!", err)
			}

			batchDocs = make([]amibaModels.DataDoc, 0)
		}
	}
	if len(batchDocs) > 0 {
		_, err := s.repo.Table("suite_amiba_data_docs").Insert(&batchDocs)
		if err != nil {
			glog.CheckAndPrintError("insert into docs error!", err)
		}
	}
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,插入单据头数据:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Name, period.Name, len(dataDocs), time.Now().Sub(fm_time).Seconds())

	//分批插入行
	fm_time = time.Now()
	batchDocLines := make([]amibaModels.DataDocLine, 0)
	count = 0
	for _, item := range dataDocLines {
		count++
		batchDocLines = append(batchDocLines, item)
		if count > 1000 {
			_, err := s.repo.Table("suite_amiba_data_doc_lines").Insert(&batchDocLines)
			if err != nil {
				glog.CheckAndPrintError("insert into docs error!", err)
			}

			batchDocLines = make([]amibaModels.DataDocLine, 0)
		}
	}
	if len(batchDocLines) > 0 {
		_, err := s.repo.Table("suite_amiba_data_doc_lines").Insert(&batchDocLines)
		if err != nil {
			glog.CheckAndPrintError("insert into docs error!", err)
		}
	}
	glog.Printf("企业:%v,核算:%v,模型:%v,期间:%v,插入单据行数据:%v条,time:%v Seconds", ent.Name, purpose.Name, m.Name, period.Name, len(dataDocLines), time.Now().Sub(fm_time).Seconds())

}
