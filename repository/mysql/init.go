/**
 @author: robert
 @date: 2021/3/4
**/
package mysql

import (
	"Garyen-go/pkg/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	db  *gorm.DB
	err error
)

// init mysql conn
func init() {
	user := setting.MySQLUser
	password := setting.MySQLPass
	host := setting.MySQLHost
	port := setting.MySQLPort
	dbName := setting.MySQLDatabase

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, password, host, port, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(15)
}
