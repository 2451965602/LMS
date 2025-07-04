// Code generated by hertz generator.

package user

import (
	"context"

	mw "github.com/2451965602/LMS/biz/middleware"
	"github.com/2451965602/LMS/biz/pack"
	"github.com/2451965602/LMS/biz/service"
	metainfoContext "github.com/2451965602/LMS/pkg/base/context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/model/user"
)

// Register .
// @router /user/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.RegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(user.RegisterResponse)

	userId, err := service.NewUserService(ctx, c).Register(ctx, req.Username, req.Password, req.Phone)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.UserID = userId

	pack.SendResponse(c, resp)
}

// Login .
// @router /user/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(user.LoginResponse)

	info, err := service.NewUserService(ctx, c).Login(ctx, req.Username, req.Password)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	mw.AccessTokenJwtMiddleware.LoginHandler(ctx, c)
	mw.RefreshTokenJwtMiddleware.LoginHandler(ctx, c)

	resp.Base = pack.BuildBaseResp(nil)
	resp.Data = pack.BuildUserResp(info)

	c.Header("Access-Token", c.GetString("Access-Token"))
	c.Header("Refresh-Token", c.GetString("Refresh-Token"))

	pack.SendResponse(c, resp)
}

// UpdateUser .
// @router /user/update [PUT]
func UpdateUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UpdateUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(user.UpdateUserResponse)

	info, err := service.NewUserService(ctx, c).UpdateUser(ctx, req)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Data = pack.BuildUserResp(info)

	pack.SendResponse(c, resp)
}

// DeleteUser .
// @router /user/delete [DELETE]
func DeleteUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DeleteUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(user.DeleteUserResponse)

	err = service.NewUserService(ctx, c).DeleteUser(ctx, req.Username)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(nil)
	pack.SendResponse(c, resp)
}

// GetUserInfo .
// @router /user/info [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserInfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(user.UserInfoResponse)

	info, err := service.NewUserService(ctx, c).GetUserById(ctx, req.UserID)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Data = pack.BuildUserResp(info)
	pack.SendResponse(c, resp)
}

// RefreshToken .
// @router /user/refresh [POST]
func RefreshToken(ctx context.Context, c *app.RequestContext) {
	resp := new(user.RefreshTokenResponse)
	userid, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	mw.GenerateAccessToken(c, userid)

	resp.Base = pack.BuildBaseResp(nil)
	pack.SendResponse(c, resp)
}
