/**
 @author: robert
 @date: 2021/8/19
**/
package login

import (
	rs "Garyen-go/service/base"
	"Garyen-go/service/login"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var f = func(resp interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, resp)
}

func UserLogin(c *gin.Context) {
	form := new(login.UserForm)

	if err := c.ShouldBindJSON(form); err != nil {
		log.Printf("request err, err: %s", err)
		f(rs.NewResp().Failed(1, "request param error"), c)
		return
	}

	loginService, err := login.NewCoreAccountService(form)
	if err != nil {
		f(rs.NewResp().Failed(1, err.Error()), c)
		return
	}

	token, err := loginService.Login()
	if err != nil {
		f(rs.NewResp().Failed(1, err.Error()), c)
		return
	}

	f(rs.NewRespData(token).Success(), c)
}
