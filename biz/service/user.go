package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/user"
	contextLogin "github.com/2451965602/LMS/pkg/base/context"
	"github.com/2451965602/LMS/pkg/errno"
)

// UserService 用于管理用户相关的业务逻辑，封装了用户注册、登录、更新、删除等操作。
type UserService struct {
	ctx context.Context     // 上下文，用于传递请求相关的元数据
	c   *app.RequestContext // Hertz框架的请求上下文，用于处理HTTP请求
}

// NewUserService 创建一个新的UserService实例，初始化上下文和请求上下文。
func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{
		ctx: ctx,
		c:   c,
	}
}

// Register 注册新用户
// 参数：
//   - ctx: 上下文
//   - username: 用户名
//   - password: 密码
//
// 返回值：
//   - int64: 注册成功用户的ID
//   - error: 错误信息，如果用户已存在或注册失败会返回错误
func (s *UserService) Register(ctx context.Context, username, password, phone string) (int64, error) {

	err := RegisterCheck(username, phone)
	if err != nil {
		return 0, err
	}

	exit, err := db.IsUserExist(ctx, username) // 检查用户是否已存在
	if err != nil {
		return 0, err
	}
	if exit {
		return 0, errno.Errorf(errno.ServiceUserExist, "user already exist") // 如果用户已存在，返回错误
	}

	id, err := db.RegisterUser(ctx, username, password) // 调用数据库操作函数注册用户
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Login 用户登录
// 参数：
//   - ctx: 上下文
//   - username: 用户名
//   - password: 密码
//
// 返回值：
//   - *db.User: 登录成功返回用户信息
//   - error: 错误信息，如果登录失败会返回错误
func (s *UserService) Login(ctx context.Context, username, password string) (*db.User, error) {
	info, err := db.LoginUser(ctx, username, password) // 调用数据库操作函数进行用户登录验证
	if err != nil {
		return nil, err
	}
	return info, nil
}

// UpdateUser 更新用户信息
// 参数：
//   - ctx: 上下文
//   - req: 更新用户请求，包含要更新的用户信息
//
// 返回值：
//   - *db.User: 更新成功返回用户信息
//   - error: 错误信息，如果更新失败会返回错误
func (s *UserService) UpdateUser(ctx context.Context, req user.UpdateUserRequest) (*db.User, error) {
	userId, err := contextLogin.GetLoginData(ctx) // 从上下文中获取当前登录用户ID
	if err != nil {
		return nil, err
	}

	info, err := db.UpdateUser(ctx, userId, req) // 调用数据库操作函数更新用户信息
	if err != nil {
		return nil, err
	}
	return info, nil
}

// DeleteUser 删除用户
// 参数：
//   - ctx: 上下文
//   - username: 要删除的用户名
//
// 返回值：
//   - error: 错误信息，如果删除失败会返回错误
func (s *UserService) DeleteUser(ctx context.Context, username string) error {
	currentUserID, err := contextLogin.GetLoginData(ctx) // 获取当前登录用户ID
	if err != nil {
		return err
	}

	err = db.DeleteUser(ctx, currentUserID, username) // 调用数据库操作函数删除用户
	if err != nil {
		return err
	}
	return nil
}

// GetUserById 根据用户ID获取用户信息
// 参数：
//   - ctx: 上下文
//   - userId: 用户ID
//
// 返回值：
//   - *db.User: 获取成功返回用户信息
//   - error: 错误信息，如果获取失败会返回错误
func (s *UserService) GetUserById(ctx context.Context, userId int64) (*db.User, error) {
	info, err := db.GetUserById(ctx, userId) // 调用数据库操作函数根据ID获取用户信息
	if err != nil {
		return nil, err
	}
	return info, nil
}

// GetUserByName 根据用户名获取用户信息
// 参数：
//   - ctx: 上下文
//   - username: 用户名
//
// 返回值：
//   - *db.User: 获取成功返回用户信息
//   - error: 错误信息，如果获取失败会返回错误
func (s *UserService) GetUserByName(ctx context.Context, username string) (*db.User, error) {
	info, err := db.GetUserByName(ctx, username) // 调用数据库操作函数根据用户名获取用户信息
	if err != nil {
		return nil, err
	}
	return info, nil
}

// AdminUpdateUser 管理员更新用户信息
// 参数：
//   - ctx: 上下文
//   - req: 管理员更新用户请求，包含要更新的用户信息
//
// 返回值：
//   - *db.User: 更新成功返回用户信息
//   - error: 错误信息，如果更新失败会返回错误
func (s *UserService) AdminUpdateUser(ctx context.Context, req user.AdminUpdateUserRequest) (*db.User, error) {
	currentUserID, err := contextLogin.GetLoginData(ctx) // 获取当前登录用户ID
	if err != nil {
		return nil, err
	}
	ok, err := db.IsPermission(ctx, currentUserID, "admin") // 检查当前用户是否有管理员权限
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errno.Errorf(errno.ServicePermissionDenied, "permission denied") // 如果没有权限，返回错误
	}

	info, err := db.AdminUpdateUser(ctx, req) // 调用数据库操作函数进行管理员更新用户操作
	if err != nil {
		return nil, err
	}
	return info, nil
}

// AdminDeleteUser 管理员删除用户
// 参数：
//   - ctx: 上下文
//   - req: 管理员删除用户请求，包含要删除的用户ID
//
// 返回值：
//   - error: 错误信息，如果删除失败会返回错误
func (s *UserService) AdminDeleteUser(ctx context.Context, req user.AdminDeleteUserRequest) error {
	currentUserID, err := contextLogin.GetLoginData(ctx) // 获取当前登录用户ID
	if err != nil {
		return err
	}
	ok, err := db.IsPermission(ctx, currentUserID, "admin") // 检查当前用户是否有管理员权限
	if err != nil {
		return err
	}
	if !ok {
		return errno.Errorf(errno.ServicePermissionDenied, "permission denied") // 如果没有权限，返回错误
	}

	err = db.AdminDeleteUser(ctx, req.UserID) // 调用数据库操作函数进行管理员删除用户操作
	if err != nil {
		return err
	}
	return nil
}
