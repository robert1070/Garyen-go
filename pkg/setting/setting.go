package setting

import "time"

const (
	HttpPort        = 8040
	ReadTimeoutSec  = 60
	WriteTimeoutSec = 60
	RunMode         = "debug"
	JwtSecret       = "fffffdsafsaddfdg"
)

const (
	MySQLUser     = "root"
	MySQLPass     = "root"
	MySQLPort     = 3306
	MySQLHost     = "127.0.0.1"
	MySQLDatabase = "test1"
)

var (
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
)

func init() {
	ReadTimeout = time.Duration(ReadTimeoutSec) * time.Second
	WriteTimeout = time.Duration(WriteTimeoutSec) * time.Second
}
