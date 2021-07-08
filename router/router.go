package router

import (
	"Garyen-go/middleware/cors"
	"Garyen-go/middleware/user"
	"Garyen-go/pkg/setting"
	v1 "Garyen-go/router/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Cors())
	r.Use(user.User())
	gin.SetMode(setting.RunMode)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/heart", v1.Ping)
	}
	return r
}
