package config

const (
	DbName    = "mysql"
	DbSchema  = "default"
	DbURL     = "root:root@tcp(127.0.0.1:3306)/default?parseTime=true"
	RedisURL  = "localhost:6379"
	RedisPass = ""
	RedisDB   = 0
	AuthID = "amit@auth.com"
	AuthPass = "amit123"
)

var (
	SecretKey = []byte("global")
)
