package constants

import "time"

const (
	AccessTokenTTL  = time.Hour * 24     //定义了Access Token的有效期
	RefreshTokenTTL = time.Hour * 24 * 7 //定义了Refresh Token的有效期
	IdentityKey     = "user_id"          //在JWT的Claims中，使用此键名存储用户ID
)
