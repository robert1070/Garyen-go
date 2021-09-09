/**
 @author: robert
 @date: 2021/9/9
**/
package lib

import (
	"errors"
	"time"

	"Garyen-go/pkg/setting"
	"github.com/dgrijalva/jwt-go"
)

// Generate jwt
func JwtAuth(user string, role string) (string, error) {
	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	st, err := token.SignedString([]byte(setting.JwtSecret))
	if err != nil {
		return "", errors.New("jwt Generate Failure")
	}

	return st, nil
}
