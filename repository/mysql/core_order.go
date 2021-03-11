package mysql

import (
	"Garyen-go/model"
	"log"
	"time"
)

type CoreOrderRepo struct {
	TableName string `json:"table_name"`
}

func (c *CoreOrderRepo) CreateCoreOrder(order *model.CoreOrder) (insertId int64, err error) {
	order.GmtCreate = time.Now().Unix()
	order.GmtModified = time.Now().Unix()
	err = db.Table(c.TableName).Create(order).Error
	if err == nil {
		insertId = order.ID
	}

	if err != nil {
		log.Printf("create core order falied, err: %v", err)
	}
	return
}
