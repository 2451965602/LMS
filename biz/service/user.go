package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/user"
	contextLogin "github.com/2451965602/LMS/pkg/base/context"
	"github.com/2451965602/LMS/pkg/errno"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{
		ctx: ctx,
		c:   c,
	}
}

func (s *UserService) Register(ctx context.Context, username, password string) (int64, error) {
	exit, err := db.IsUserExist(ctx, username)
	if err != nil {
		return 0, err
	}
	if exit {
		return 0, errno.Errorf(errno.ServiceUserExist, "user already exist")
	}

	id, err := db.RegisterUser(ctx, username, password)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*db.User, error) {
	info, err := db.LoginUser(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req user.UpdateUserRequest) (*db.User, error) {
	userId, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, err
	}

	info, err := db.UpdateUser(ctx, userId, req)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *UserService) DeleteUser(ctx context.Context, username string) error {
	currentUserID, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return err
	}

	err = db.DeleteUser(ctx, currentUserID, username)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserById(ctx context.Context, userId int64) (*db.User, error) {
	info, err := db.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *UserService) GetUserByName(ctx context.Context, username string) (*db.User, error) {
	info, err := db.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *UserService) AdminUpdateUser(ctx context.Context, req user.AdminUpdateUserRequest) (*db.User, error) {
	currentUserID, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, err
	}
	ok, err := db.IsPermission(ctx, currentUserID, "admin")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errno.Errorf(errno.ServicePermissionDenied, "permission denied")
	}

	info, err := db.AdminUpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *UserService) AdminDeleteUser(ctx context.Context, req user.AdminDeleteUserRequest) error {
	currentUserID, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return err
	}
	ok, err := db.IsPermission(ctx, currentUserID, "admin")
	if err != nil {
		return err
	}
	if !ok {
		return errno.Errorf(errno.ServicePermissionDenied, "permission denied")
	}

	err = db.AdminDeleteUser(ctx, req.UserID)
	if err != nil {
		return err
	}
	return nil
}
