package auth

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	mw "github.com/2451965602/LMS/biz/middleware"
	"github.com/2451965602/LMS/biz/pack"
	metainfoContext "github.com/2451965602/LMS/pkg/base/context"
	"github.com/2451965602/LMS/pkg/errno"
)

func AccessTokenAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		valid, userID := mw.IsAccessTokenAvailable(ctx, c)
		if !valid {
			pack.SendFailResponse(c, errno.Errorf(errno.AuthAccessExpiredCode, "access token expired"))
			c.Abort()
			return
		}
		ctx = metainfoContext.WithLoginData(ctx, userID)
		c.Next(ctx)
	}
}

func RefreshTokenAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		valid, userID := mw.IsRefreshTokenAvailable(ctx, c)
		if !valid {
			pack.SendFailResponse(c, errno.Errorf(errno.AuthRefreshExpiredCode, "refresh token expired"))
			c.Abort()
			return
		}
		ctx = metainfoContext.WithLoginData(ctx, userID)
		c.Next(ctx)
	}
}
