/**
 @author: robert
 @date: 2021/3/12
**/
package operator

import (
	"Garyen-go/proto"
)

var (
	_ Operator = &Pager{}
)

type Operator interface {
	Handler(args ...interface{}) error
}

type Pager struct{}

func (p *Pager) Handler(params ...interface{}) error {
	if len(params) == 0 {
		return nil
	}

	for _, param := range params {
		var page *proto.PageRequest
		switch val := param.(type) {
		case *proto.PageRequest:
			page = val
		case *proto.IDAndPageRequest:
			page = &val.PageRequest
		default:
			continue
		}

		if page.Page <= 0 {
			page.Page = 1
		}

		if page.PageSize <= 0 {
			page.PageSize = 20
		}

		page.Offset = (page.Page - 1) * page.PageSize
	}

	return nil
}
