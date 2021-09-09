/**
 @author: robert
 @date: 2021/3/4
**/
package dms

import (
	dms2 "Garyen-go/model/dms"
	"Garyen-go/repository/mysql/dms"
)

func NewCoreOrderService() *Service {
	return &Service{
		mysql: &dms.CoreOrderRepo{
			TableName: "core_sql_order",
		},
	}
}

type Service struct {
	mysql *dms.CoreOrderRepo
}

func (srv *Service) Create(order *dms2.CoreSQLOrder) bool {
	_, err := srv.mysql.Add(order)
	if err != nil {
		return false
	}
	return true
}
