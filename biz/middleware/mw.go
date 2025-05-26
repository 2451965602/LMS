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
	AccessTokenJwtMiddleware  *jwt.HertzJWTMiddleware // Access Token的JWT中间件实例
	RefreshTokenJwtMiddleware *jwt.HertzJWTMiddleware // Refresh Token的JWT中间件实例
)

// initJWTCommonConfig 初始化JWT中间件的通用配置
func initJWTCommonConfig() *jwt.HertzJWTMiddleware {
	return &jwt.HertzJWTMiddleware{
		Realm:                       "LMS",                 // JWT的Realm名称
		WithoutDefaultTokenHeadName: true,                  // 不使用默认的Token头名称
		IdentityKey:                 constants.IdentityKey, // 用户身份标识的键名
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			pack.SendFailResponse(c, errno.AuthInvalid) // 未授权时的响应处理
		},
	}
}

// AccessTokenJwt 初始化Access Token的JWT中间件
func AccessTokenJwt() {
	config := initJWTCommonConfig()
	config.Key = []byte("AccessToken_key")       // Access Token的密钥
	config.Timeout = constants.AccessTokenTTL    // Access Token的有效期
	config.TokenLookup = "header: Authorization" // 从请求头的Authorization字段查找Token

	config.PayloadFunc = func(data interface{}) jwt.MapClaims {
		return jwt.MapClaims{
			config.IdentityKey: data,     // 用户身份标识
			"token_type":       "access", // Token类型为Access Token
		}
	}

	config.IdentityHandler = func(ctx context.Context, c *app.RequestContext) interface{} {
		claims := jwt.ExtractClaims(ctx, c) // 提取JWT中的Claims
		return claims[config.IdentityKey]   // 返回用户身份标识
	}

	config.LoginResponse = func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
		c.Set("Access-Token", token) // 将Access Token设置到上下文中
	}

	config.Authenticator = func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
		var loginStruct user.LoginRequest
		if err := c.BindAndValidate(&loginStruct); err != nil { // 绑定并验证登录请求
			return nil, err
		}
		users, err := db.LoginUser(ctx, loginStruct.Username, loginStruct.Password) // 验证用户登录
		if err != nil {
			return nil, err
		}
		ctx = metainfoContext.WithLoginData(ctx, users.ID) // 将用户ID设置到上下文中
		return users.ID, nil
	}

	config.Authorizator = func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
		claims := jwt.ExtractClaims(ctx, c)                                              // 提取JWT中的Claims
		if tokenType, ok := claims["token_type"].(string); ok && tokenType == "access" { // 检查Token类型是否为Access Token
			return true
		}
		return false
	}

	var err error
	AccessTokenJwtMiddleware, err = jwt.New(config) // 创建JWT中间件实例
	if err != nil {
		hlog.Fatal("AccessToken JWT Error:" + err.Error()) // 如果初始化失败，记录错误日志并退出
	}
}

// RefreshTokenJwt 初始化Refresh Token的JWT中间件
func RefreshTokenJwt() {
	config := initJWTCommonConfig()
	config.Key = []byte("refresh_token_key")     // Refresh Token的密钥
	config.Timeout = constants.RefreshTokenTTL   // Refresh Token的有效期
	config.TokenLookup = "header: Refresh-Token" // 从请求头的Refresh-Token字段查找Token

	config.PayloadFunc = func(data interface{}) jwt.MapClaims {
		return jwt.MapClaims{
			config.IdentityKey: data,      // 用户身份标识
			"token_type":       "refresh", // Token类型为Refresh Token
		}
	}

	config.IdentityHandler = func(ctx context.Context, c *app.RequestContext) interface{} {
		claims := jwt.ExtractClaims(ctx, c) // 提取JWT中的Claims
		return claims[config.IdentityKey]   // 返回用户身份标识
	}

	config.LoginResponse = func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
		c.Set("Refresh-Token", token) // 将Refresh Token设置到上下文中
	}

	config.Authenticator = func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
		var loginStruct user.LoginRequest
		if err := c.BindAndValidate(&loginStruct); err != nil { // 绑定并验证登录请求
			return nil, err
		}
		users, err := db.LoginUser(ctx, loginStruct.Username, loginStruct.Password) // 验证用户登录
		if err != nil {
			return nil, err
		}
		ctx = metainfoContext.WithLoginData(ctx, users.ID) // 将用户ID设置到上下文中
		return users.ID, nil
	}

	config.Authorizator = func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
		claims := jwt.ExtractClaims(ctx, c)                                               // 提取JWT中的Claims
		if tokenType, ok := claims["token_type"].(string); ok && tokenType == "refresh" { // 检查Token类型是否为Refresh Token
			return true
		}
		return false
	}

	var err error
	RefreshTokenJwtMiddleware, err = jwt.New(config) // 创建JWT中间件实例
	if err != nil {
		hlog.Fatal("RefreshToken JWT Error:" + err.Error()) // 如果初始化失败，记录错误日志并退出
	}
}

// GenerateAccessToken 生成新的Access Token
func GenerateAccessToken(c *app.RequestContext, userId int64) {
	tokenString, _, _ := AccessTokenJwtMiddleware.TokenGenerator(userId) // 调用Token生成器生成新的Access Token
	c.Header("New-Access-Token", tokenString)                            // 将新的Access Token设置到响应头中
}

// IsAccessTokenAvailable 检查Access Token是否有效
func IsAccessTokenAvailable(ctx context.Context, c *app.RequestContext) (bool, int64) {
	claims, err := AccessTokenJwtMiddleware.GetClaimsFromJWT(ctx, c) // 从JWT中提取Claims
	if err != nil {
		return false, 0
	}

	if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "access" { // 检查Token类型是否为Access Token
		return false, 0
	}

	switch v := claims["exp"].(type) {
	case nil:
		return false, 0
	case float64:
		if int64(v) < AccessTokenJwtMiddleware.TimeFunc().Unix() { // 检查Token是否过期
			return false, 0
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false, 0
		}
		if n < AccessTokenJwtMiddleware.TimeFunc().Unix() { // 检查Token是否过期
			return false, 0
		}
	default:
		return false, 0
	}

	c.Set("JWT_PAYLOAD", claims)                                 // 将JWT的Claims设置到上下文中
	identity := AccessTokenJwtMiddleware.IdentityHandler(ctx, c) // 获取用户身份标识

	var userID int64
	if id, ok := claims[AccessTokenJwtMiddleware.IdentityKey].(float64); ok { // 提取用户ID
		userID = int64(id)
	} else {
		return false, 0
	}

	if identity != nil {
		c.Set(AccessTokenJwtMiddleware.IdentityKey, identity) // 将用户身份标识设置到上下文中
	}

	isValid := AccessTokenJwtMiddleware.Authorizator(identity, ctx, c) // 检查Token是否有效
	return isValid, userID
}

// IsRefreshTokenAvailable 检查Refresh Token是否有效
func IsRefreshTokenAvailable(ctx context.Context, c *app.RequestContext) (bool, int64) {
	claims, err := RefreshTokenJwtMiddleware.GetClaimsFromJWT(ctx, c) // 从JWT中提取Claims
	if err != nil {
		return false, 0
	}

	if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "refresh" { // 检查Token类型是否为Refresh Token
		return false, 0
	}

	switch v := claims["exp"].(type) {
	case nil:
		return false, 0
	case float64:
		if int64(v) < RefreshTokenJwtMiddleware.TimeFunc().Unix() { // 检查Token是否过期
			return false, 0
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false, 0
		}
		if n < RefreshTokenJwtMiddleware.TimeFunc().Unix() { // 检查Token是否过期
			return false, 0
		}
	default:
		return false, 0
	}

	c.Set("JWT_PAYLOAD", claims)                                  // 将JWT的Claims设置到上下文中
	identity := RefreshTokenJwtMiddleware.IdentityHandler(ctx, c) // 获取用户身份标识

	var userID int64
	if id, ok := claims[RefreshTokenJwtMiddleware.IdentityKey].(float64); ok { // 提取用户ID
		userID = int64(id)
	} else {
		return false, 0
	}

	if identity != nil {
		c.Set(RefreshTokenJwtMiddleware.IdentityKey, identity) // 将用户身份标识设置到上下文中
	}

	isValid := RefreshTokenJwtMiddleware.Authorizator(identity, ctx, c) // 检查Token是否有效
	return isValid, userID
}

// Init 初始化JWT中间件
func Init() {
	AccessTokenJwt()  // 初始化Access Token的JWT中间件
	RefreshTokenJwt() // 初始化Refresh Token的JWT中间件

	if err := AccessTokenJwtMiddleware.MiddlewareInit(); err != nil { // 初始化Access Token的JWT中间件
		hlog.Fatal("AccessTokenJwtMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	if err := RefreshTokenJwtMiddleware.MiddlewareInit(); err != nil { // 初始化Refresh Token的JWT中间件
		hlog.Fatal("RefreshTokenJwtMiddleware.MiddlewareInit() Error:" + err.Error())
	}
}
