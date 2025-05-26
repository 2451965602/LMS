package context

import (
	"context"
	"strconv"

	"github.com/2451965602/LMS/pkg/constants"
	"github.com/2451965602/LMS/pkg/errno"
)

// WithLoginData 将用户ID存储到上下文中
// 参数：
//   - ctx: 原始上下文
//   - uid: 用户ID
//
// 返回值：
//   - context.Context: 包含用户ID的上下文
func WithLoginData(ctx context.Context, uid int64) context.Context {
	return newContext(ctx, constants.LoginDataKey, strconv.FormatInt(uid, 10)) // 将用户ID存储到上下文中
}

// GetLoginData 从上下文中获取用户ID
// 参数：
//   - ctx: 包含用户ID的上下文
//
// 返回值：
//   - int64: 用户ID
//   - error: 错误信息，如果获取失败会返回错误
func GetLoginData(ctx context.Context) (int64, error) {
	user, ok := fromContext(ctx, constants.LoginDataKey) // 从上下文中获取用户ID
	if !ok {
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get header in context") // 如果未找到用户ID，返回错误
	}

	value, err := strconv.ParseInt(user, 10, 64) // 将用户ID从字符串转换为整数
	if err != nil {
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to get header in context when parse loginData") // 如果转换失败，返回错误
	}
	return value, nil
}
