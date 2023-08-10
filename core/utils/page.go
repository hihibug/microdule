package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pages struct {
	Page int `json:"page" binding:"required"`
	Size int `json:"page_size" binding:"required"`
}

// SearchPageParams api分页条件查询及排序结构体
type SearchPageParams struct {
	OrderKey string `json:"order_key"` // 排序
	Desc     bool   `json:"desc"`      // 排序方式:升序false(默认)|降序true
	Pages           //分页
}

// PageQuery 分页请求
type PageQuery struct {
	Model       *gorm.DB
	PageStruct  SearchPageParams
	OrderKeys   []string
	OrderPrefix string
}

// BuffPageFmt 格式化分页数据
func BuffPageFmt(orderKey string, page, pageSize int, orderDesc bool) SearchPageParams {
	sp := SearchPageParams{
		OrderKey: orderKey,
		Desc:     false,
		Pages: Pages{
			Page: page,
			Size: pageSize,
		},
	}
	if orderDesc {
		sp.Desc = true
	}
	return sp
}

// GeneratePaginationFromRequest 生成分页
func GeneratePaginationFromRequest(c *gin.Context) (SearchPageParams, error) {
	var pagination SearchPageParams
	if err := c.ShouldBind(&pagination); err != nil {
		return SearchPageParams{}, err
	}

	// 校验参数
	if pagination.Size <= 0 {
		pagination.Size = 10
	}

	if pagination.Page <= 1 {
		pagination.Page = 1
	}

	return pagination, nil
}

// OrderByString 排序字段定义
func OrderByString(key []string, order string, desc bool) (string, error) {
	// 设置有效排序key 防止sql注入
	orderMap := make(map[string]bool, len(key))
	for _, v := range key {
		orderMap[v] = true
	}
	if !orderMap[order] {
		return "", fmt.Errorf("排序字段错误")
	}
	if desc {
		order = order + " desc"
	}
	return order, nil
}

// SearchPageFmt sql Page 封装
func SearchPageFmt(model *gorm.DB, q SearchPageParams, orderKeys []string, orderPrefix string) (total int64, pageNum int64, db *gorm.DB, err error) {
	page := &PageQuery{
		Model:       model,
		PageStruct:  q,
		OrderKeys:   orderKeys,
		OrderPrefix: orderPrefix,
	}

	total, pageNum, err = page.PageTotal()
	if err != nil {
		return 0, 0, page.Model, err
	}

	db, err = page.PageOrderBy()
	if err != nil {
		return total, pageNum, db, err
	}

	return total, pageNum, db, nil
}

// GroupSearchPageFmt group Page 封装
func GroupSearchPageFmt(model *gorm.DB, q SearchPageParams, orderKeys []string, orderPrefix string) (total int64, pageNum int64, db *gorm.DB, err error) {
	page := &PageQuery{
		Model:       model,
		PageStruct:  q,
		OrderKeys:   orderKeys,
		OrderPrefix: orderPrefix,
	}

	total, pageNum, err = page.GroupPageTotal()
	if err != nil {
		return 0, 0, page.Model, err
	}

	db, err = page.PageOrderBy()
	if err != nil {
		return total, pageNum, db, err
	}

	return total, pageNum, db, nil
}

// ParallelSearchPageFmt 并行分页
func ParallelSearchPageFmt(model *gorm.DB, q SearchPageParams, orderKeys []string, orderPrefix string, result interface{}) (total int64, pageNum int64, err error) {
	page := &PageQuery{
		Model:       model,
		PageStruct:  q,
		OrderKeys:   orderKeys,
		OrderPrefix: orderPrefix,
	}

	f0 := func() error {
		total, pageNum, err = page.PageTotal()
		if err != nil {
			return err
		}
		return nil
	}

	f1 := func() error {
		dbs, err := page.PageOrderBy()
		if err != nil {
			return err
		}
		dbs.Find(result)
		return nil
	}

	err = GoPanic(f0, f1)
	if err != nil {
		return 0, 0, err
	}

	return total, pageNum, nil
}

// ParallelGroupSearchPageFmt 并行分页
func ParallelGroupSearchPageFmt(model *gorm.DB, q SearchPageParams, orderKeys []string, orderPrefix string, result interface{}) (total int64, pageNum int64, err error) {
	page := &PageQuery{
		Model:       model,
		PageStruct:  q,
		OrderKeys:   orderKeys,
		OrderPrefix: orderPrefix,
	}

	f0 := func() error {
		total, pageNum, err = page.GroupPageTotal()
		if err != nil {
			return err
		}
		return nil
	}

	f1 := func() error {
		dbs, err := page.PageOrderBy()
		if err != nil {
			return err
		}
		dbs.Find(result)
		return nil
	}

	err = GoPanic(f0, f1)
	if err != nil {
		return 0, 0, err
	}

	return total, pageNum, nil
}

// PageTotal 统计总数 获取当前页数
func (p *PageQuery) PageTotal() (total, pageNum int64, err error) {
	err = p.Model.Session(&gorm.Session{}).Count(&total).Error
	pageNum = total / int64(p.PageStruct.Size)
	if total%int64(p.PageStruct.Size) != 0 {
		pageNum++
	}
	return
}

// GroupPageTotal 统计总数 获取当前页数
func (p *PageQuery) GroupPageTotal() (total, pageNum int64, err error) {
	//分页
	count := p.Model.Session(&gorm.Session{}).Count(&total)
	total = count.RowsAffected
	err = count.Error

	pageNum = total / int64(p.PageStruct.Size)
	if total%int64(p.PageStruct.Size) != 0 {
		pageNum++
	}
	return
}

// PageOrderBy 组装排序
func (p *PageQuery) PageOrderBy() (db *gorm.DB, err error) {
	//排序
	model := p.Model.Session(&gorm.Session{})
	if p.PageStruct.OrderKey != "" {
		orderKey, err := OrderByString(p.OrderKeys[:], p.PageStruct.OrderKey, p.PageStruct.Desc)
		if err != nil {
			return p.Model, err
		}
		if p.OrderPrefix != "" {
			p.Model = model.Order(p.OrderPrefix + orderKey)
		} else {
			p.Model = model.Order(orderKey)
		}
	}
	db = model.Limit(p.PageStruct.Size).Offset((p.PageStruct.Page - 1) * p.PageStruct.Size)
	return
}

// PageFmt Page 格式化输出
func PageFmt(data interface{}, total, pageNum int64) interface{} {
	list := make(map[string]interface{})
	list["lists"] = data
	list["total"] = total
	list["page"] = pageNum

	return list
}
