package models

type RequestQuery struct {
	Search   string `json:"search" form:"search"`
	PageNo   int64  `json:"pageNo" form:"pageNo"`
	PageSize int64  `json:"pageSize" form:"pageSize"`
}
