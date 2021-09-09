package v1

import (
	rs "Garyen-go/service/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, rs.NewResp().Success())
}
