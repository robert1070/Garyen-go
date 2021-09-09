/**
 @author: robert
 @date: 2021/8/19
**/
package login

import (
	"Garyen-go/lib"
	u "Garyen-go/repository/mysql/user"
	"errors"
)

type LType int8

const (
	g LType = iota + 1
	l
)

type UserForm struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	LoginType LType  `json:"login_type"`
}

func NewCoreAccountService(form *UserForm) (UserLogin, error) {
	var loginService UserLogin
	switch form.LoginType {
	case g:
		loginService = &general{
			user: form,
			mysql: &u.CoreAccountRepo{
				TableName: "core_account",
			},
		}
	case l:
		loginService = &ldap{
			user: form,
			mysql: &u.CoreAccountRepo{
				TableName: "core_account",
			},
		}
	default:
		return nil, errors.New("not support login type")
	}

	return loginService, nil
}

type UserLogin interface {
	Login() (map[string]interface{}, error)
}

type general struct {
	user  *UserForm
	mysql *u.CoreAccountRepo
}

func (g *general) Login() (map[string]interface{}, error) {
	account, err := g.mysql.Get(g.user.UserName)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("user not exists")
	}

	if !lib.DjangoCheckPassword(g.user.Password, account.Password) {
		return nil, errors.New("incorrect password")
	}

	token, err := lib.JwtAuth(g.user.Password, account.Rule)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

type ldap struct {
	user  *UserForm
	mysql *u.CoreAccountRepo
}

func (l *ldap) Login() (map[string]interface{}, error) {
	//todo
	return nil, nil
}
