package jwt

import (
	"Garyen-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type jwtPayload struct {
	Exp  int64  `json:"exp"`
	Name string `json:"name"`
	Role string `json:"role"`
}

const (
	USER = "user"
)

func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Status(http.StatusUnauthorized)
			return
		}

		arr := strings.Split(auth, ".")
		if len(arr) != 3 {
			c.Status(http.StatusUnauthorized)
			return
		}

		jpStr, err := utils.Base64Decode(arr[1])
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		var jwt jwtPayload
		if err := utils.JSONDecodeFromString(jpStr, &jwt); err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Set(USER, jwt)
		c.Next()
	}
}
