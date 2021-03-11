package v1

import (
	"Garyen-go/service/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, vo.HandlerRespVo(vo.BaseRespType).Success())
}