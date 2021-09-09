/**
 @author: robert
 @date: 2021/8/19
**/
package user

import (
	u "Garyen-go/model/user"
	"Garyen-go/repository/mysql"
	"time"
)

type CoreAccountRepo struct {
	TableName string `json:"table_name"`
}

func (c *CoreAccountRepo) Add() error {
	currentTime := time.Now().Unix()
	if err := mysql.DB.Table(c.TableName).Create(&u.CoreAccount{
		Username:    "",
		Password:    "",
		Rule:        "develop",
		RealName:    "",
		Email:       "",
		State:       1,
		GmtCreate:   currentTime,
		GmtModified: currentTime,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (c *CoreAccountRepo) Get(user string) (*u.CoreAccount, error) {
	var account u.CoreAccount
	if err := mysql.DB.Table(c.TableName).Where("username = ?", user).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *CoreAccountRepo) Update() bool {
	return false
}
