package lib

import (
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
	"strconv"
	"strings"
)

func DjangoEncrypt(password string, sl string) string {
	pwd := []byte(password)
	salt := []byte(sl)
	iterations := 120000
	digest := sha256.New
	dk := pbkdf2.Key(pwd, salt, iterations, 32, digest)
	str := base64.StdEncoding.EncodeToString(dk)
	return "pbkdf2_sha256" + "$" + strconv.FormatInt(int64(iterations), 10) + "$" + string(salt) + "$" + str
}

func DjangoCheckPassword(inputPwd, RealPwd string) bool {
	sl := strings.Split(RealPwd, "$")[2]
	encryptPwd := DjangoEncrypt(inputPwd, sl)
	if !strings.EqualFold(RealPwd, encryptPwd) {
		return false
	}

	return true
}
