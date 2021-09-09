/**
 @author: robert
 @date: 2021/7/13
**/
package proto

type IDRequest struct {
	ID int `json:"id" binding:"required"`
}

type PageRequest struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Offset   int64 `json:"offset"`
}

type IDAndPageRequest struct {
	IDRequest
	PageRequest
}
