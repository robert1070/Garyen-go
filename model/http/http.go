/**
 @author: robert
 @date: 2021/3/12
**/
package http

type IDRequest struct {
	ID int `json:"id" binding:"required"`
}

type PageRequest struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Offset   int64 `json:"-"`
}

type IDAndPageRequest struct {
	IDRequest
	PageRequest
}