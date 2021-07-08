/**
 @author: robert
 @date: 2021/3/12
**/
package operator

import "Garyen-go/model/http"

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
		var page *http.PageRequest
		switch val := param.(type) {
		case *http.PageRequest:
			page = val
		case *http.IDAndPageRequest:
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
