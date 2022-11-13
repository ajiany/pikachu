package paginator

import "gorm.io/gorm"

const DefaultPageSize = 10
const MaxPageSize = 100

type PaginationData struct {
	Count    int64       `test_helper:"count"`
	Current  int         `test_helper:"current"`
	PageSize int         `test_helper:"page_size"`
	HasNext  bool        `test_helper:"has_next"`
	Data     interface{} `test_helper:"data"`
}

func NewPaginationData(total int64, page int, pageSize int, data interface{}) *PaginationData {
	hasNext := int64(page*pageSize) < total

	return &PaginationData{
		Count:    total,
		Current:  page,
		PageSize: pageSize,
		HasNext:  hasNext,
		Data:     data,
	}
}

type Pagination struct {
	Page     int
	PageSize int
}

func NewPagination(page int, pageSize int) *Pagination {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) Query(db *gorm.DB, data interface{}, selectColumn ...string) (*PaginationData, error) {
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, err
	}

	if len(selectColumn) > 0 {
		db = db.Select(selectColumn[0])
	}

	if err := db.Limit(p.PageSize).Offset(p.GetOffset()).Scan(data).Error; err != nil {
		return nil, err
	}

	return NewPaginationData(count, p.Page, p.PageSize, data), nil
}
