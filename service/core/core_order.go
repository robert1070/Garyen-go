/**
 @author: robert
 @date: 2021/3/4
**/
package core

import (
	"Garyen-go/model"
	"Garyen-go/repository/mysql"
)

func NewCoreOrderService() *Service {
	return &Service{
		mysql: &mysql.CoreOrderRepo{
			TableName: "core_order",
		},
	}
}

type Service struct {
	mysql *mysql.CoreOrderRepo
}

func (srv *Service) Create(order *model.CoreOrder) bool {
	_, err := srv.mysql.CreateCoreOrder(order)
	if err != nil {
		return false
	}

	return true
}
