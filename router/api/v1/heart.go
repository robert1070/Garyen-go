package v1

import (
	"Garyen-go/service/dms/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, vo.HandlerRespVo(vo.BaseRespType).Success())
}