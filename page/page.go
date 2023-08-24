package page

import (
	"math"

	"gorm.io/gorm"
)

// 通用分页器
type PageInterface interface {
	// querySql 查询语句
	// countSql 统计语句
	// args 查询参数
	// dataMode 查询的数据模型
	// page 页码，从1开始
	// pageSize 每页显示数量
	List(querySql, countSql string, args []interface{}, dataMode interface{}, page, pageSize int) (result Result, err error)
}

type Result struct {
	Total    int         `json:"total"`    //数据总数
	PageNum  int         `json:"pageNum"`  //总页数
	PageSize int         `json:"pageSize"` //每页记录数
	Page     int         `json:"page"`     //当前页，从1开始
	Data     interface{} `json:"data"`     //数据
}

type Pager struct {
	db *gorm.DB
}

func NewPager(db *gorm.DB) PageInterface {
	return &Pager{db: db}
}

func (p *Pager) List(querySql, countSql string, args []interface{}, dataMode interface{}, page, pageSize int) (result Result, err error) {
	var count int64
	err = p.db.Raw(countSql, args...).Count(&count).Error
	if err != nil {
		return result, err
	}
	if count == 0 {
		return result, nil
	}
	err = p.db.Raw(querySql, args...).Limit(pageSize).Offset((page - 1) * pageSize).Scan(dataMode).Error
	if err != nil {
		return result, err
	}
	result.Total = int(count)
	result.Page = page
	result.PageSize = pageSize
	result.PageNum = int(math.Ceil(float64(count) / float64(pageSize)))
	result.Data = dataMode
	return result, nil
}
