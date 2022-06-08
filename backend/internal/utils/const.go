package utils

import "time"

// UserKey is
const (
	UserKey           = "user"
	AuthStatusKey     = "auth_status"
	CookieExpires     = 240 * time.Hour
	CsrfExpires       = 10 * time.Minute
	Host              = "localhost"
	Port              = "8085"
	Port2             = "8086"
	AuthPort          = "8081"
	FileServerPort    = "8082"
	RedisPort         = "6379"
	AvatarsPath       = "http://" + Host + ":8085" + "/avatars/"
	DefaultAvatarPath = "http://" + Host + "/avatars/default.jpeg"
	AvatarsFileDir    = "tmp/avatars/"

	RequestID = "RequestID"
)

const (
	ActionAdd = "add"
	ActionGet = "get"
)
