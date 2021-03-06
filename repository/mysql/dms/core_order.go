package dms

import (
	"Garyen-go/model/dms"
	"Garyen-go/repository/mysql"
	"log"
	"time"
)

type CoreOrderRepo struct {
	TableName string `json:"table_name"`
}

func (c *CoreOrderRepo) Add(order *dms.CoreSQLOrder) (insertId int64, err error) {
	order.GmtCreate = time.Now().Unix()
	order.GmtModified = time.Now().Unix()
	err = mysql.DB.Table(c.TableName).Create(order).Error
	if err == nil {
		insertId = order.ID
	}

	if err != nil {
		log.Printf("create core order falied, err: %v", err)
	}
	return
}
