package constants

import "time"

const (
	AccessTokenTTL  = time.Hour * 24
	RefreshTokenTTL = time.Hour * 24 * 7
	IdentityKey     = "user-id"
)
