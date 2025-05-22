package mw

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/user"
	"github.com/2451965602/LMS/biz/pack"
	metainfoContext "github.com/2451965602/LMS/pkg/base/context"
	"github.com/2451965602/LMS/pkg/constants"
	"github.com/2451965602/LMS/pkg/errno"
)

var (
	AccessTokenJwtMiddleware  *jwt.HertzJWTMiddleware
	RefreshTokenJwtMiddleware *jwt.HertzJWTMiddleware
)

func initJWTCommonConfig() *jwt.HertzJWTMiddleware {
	return &jwt.HertzJWTMiddleware{
		Realm:                       "LMS",
		WithoutDefaultTokenHeadName: true,
		IdentityKey:                 constants.IdentityKey,
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			pack.SendFailResponse(c, errno.AuthInvalid)
		},
	}
}

func AccessTokenJwt() {
	config := initJWTCommonConfig()
	config.Key = []byte("AccessToken_key")
	config.Timeout = constants.AccessTokenTTL
	config.TokenLookup = "header: Authorization"

	config.PayloadFunc = func(data interface{}) jwt.MapClaims {
		return jwt.MapClaims{
			config.IdentityKey: data,
			"token_type":       "access",
		}
	}

	config.IdentityHandler = func(ctx context.Context, c *app.RequestContext) interface{} {
		claims := jwt.ExtractClaims(ctx, c)
		return claims[config.IdentityKey]
	}

	config.LoginResponse = func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
		c.Set("Access-Token", token)
	}

	config.Authenticator = func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
		var loginStruct user.LoginRequest
		if err := c.BindAndValidate(&loginStruct); err != nil {
			return nil, err
		}
		users, err := db.LoginUser(ctx, loginStruct.Username, loginStruct.Password)
		if err != nil {
			return nil, err
		}
		ctx = metainfoContext.WithLoginData(ctx, users.ID)
		return users.ID, nil
	}

	config.Authorizator = func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
		claims := jwt.ExtractClaims(ctx, c)
		if tokenType, ok := claims["token_type"].(string); ok && tokenType == "access" {
			return true
		}
		return false
	}

	var err error
	AccessTokenJwtMiddleware, err = jwt.New(config)
	if err != nil {
		hlog.Fatal("AccessToken JWT Error:" + err.Error())
	}
}

func RefreshTokenJwt() {
	config := initJWTCommonConfig()
	config.Key = []byte("refresh_token_key")
	config.Timeout = constants.RefreshTokenTTL
	config.TokenLookup = "header: Refresh-Token"

	config.PayloadFunc = func(data interface{}) jwt.MapClaims {
		return jwt.MapClaims{
			config.IdentityKey: data,
			"token_type":       "refresh",
		}
	}

	config.IdentityHandler = func(ctx context.Context, c *app.RequestContext) interface{} {
		claims := jwt.ExtractClaims(ctx, c)
		return claims[config.IdentityKey]
	}

	config.LoginResponse = func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
		c.Set("Refresh-Token", token)
	}

	config.Authenticator = func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
		var loginStruct user.LoginRequest
		if err := c.BindAndValidate(&loginStruct); err != nil {
			return nil, err
		}
		users, err := db.LoginUser(ctx, loginStruct.Username, loginStruct.Password)
		if err != nil {
			return nil, err
		}
		ctx = metainfoContext.WithLoginData(ctx, users.ID)
		return users.ID, nil
	}

	config.Authorizator = func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
		claims := jwt.ExtractClaims(ctx, c)
		if tokenType, ok := claims["token_type"].(string); ok && tokenType == "refresh" {
			return true
		}
		return false
	}

	var err error
	RefreshTokenJwtMiddleware, err = jwt.New(config)
	if err != nil {
		hlog.Fatal("RefreshToken JWT Error:" + err.Error())
	}
}

func GenerateAccessToken(c *app.RequestContext, userId int64) {
	tokenString, _, _ := AccessTokenJwtMiddleware.TokenGenerator(userId)
	c.Header("New-Access-Token", tokenString)
}

func IsAccessTokenAvailable(ctx context.Context, c *app.RequestContext) (bool, int64) {
	claims, err := AccessTokenJwtMiddleware.GetClaimsFromJWT(ctx, c)
	if err != nil {
		return false, 0
	}

	if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "access" {
		return false, 0
	}

	switch v := claims["exp"].(type) {
	case nil:
		return false, 0
	case float64:
		if int64(v) < AccessTokenJwtMiddleware.TimeFunc().Unix() {
			return false, 0
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false, 0
		}
		if n < AccessTokenJwtMiddleware.TimeFunc().Unix() {
			return false, 0
		}
	default:
		return false, 0
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := AccessTokenJwtMiddleware.IdentityHandler(ctx, c)

	var userID int64
	if id, ok := claims[AccessTokenJwtMiddleware.IdentityKey].(float64); ok {
		userID = int64(id)
	} else {
		return false, 0
	}

	if identity != nil {
		c.Set(AccessTokenJwtMiddleware.IdentityKey, identity)
	}

	isValid := AccessTokenJwtMiddleware.Authorizator(identity, ctx, c)
	return isValid, userID
}

func IsRefreshTokenAvailable(ctx context.Context, c *app.RequestContext) (bool, int64) {
	claims, err := RefreshTokenJwtMiddleware.GetClaimsFromJWT(ctx, c)
	if err != nil {
		return false, 0
	}

	if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "refresh" {
		return false, 0
	}

	switch v := claims["exp"].(type) {
	case nil:
		return false, 0
	case float64:
		if int64(v) < RefreshTokenJwtMiddleware.TimeFunc().Unix() {
			return false, 0
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false, 0
		}
		if n < RefreshTokenJwtMiddleware.TimeFunc().Unix() {
			return false, 0
		}
	default:
		return false, 0
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := RefreshTokenJwtMiddleware.IdentityHandler(ctx, c)

	var userID int64
	if id, ok := claims[RefreshTokenJwtMiddleware.IdentityKey].(float64); ok {
		userID = int64(id)
	} else {
		return false, 0
	}

	if identity != nil {
		c.Set(RefreshTokenJwtMiddleware.IdentityKey, identity)
	}

	isValid := RefreshTokenJwtMiddleware.Authorizator(identity, ctx, c)
	return isValid, userID
}

func Init() {
	AccessTokenJwt()
	RefreshTokenJwt()

	if err := AccessTokenJwtMiddleware.MiddlewareInit(); err != nil {
		hlog.Fatal("AccessTokenJwtMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	if err := RefreshTokenJwtMiddleware.MiddlewareInit(); err != nil {
		hlog.Fatal("RefreshTokenJwtMiddleware.MiddlewareInit() Error:" + err.Error())
	}
}
