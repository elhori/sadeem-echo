package models

type Pagination struct {
	Page     int `query:"currentPage"`
	PageSize int `query:"pageSize"`
}
