package setting

import "time"

const (
	HttpPort        = 8040
	ReadTimeoutSec  = 60
	WriteTimeoutSec = 60
	RunMode         = "debug"
	JwtSecret       = ""
)

var (
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
)

func init() {
	ReadTimeout = time.Duration(ReadTimeoutSec) * time.Second
	WriteTimeout = time.Duration(WriteTimeoutSec) * time.Second
}
