package context

import (
	"context"
	"strconv"

	"github.com/2451965602/LMS/pkg/constants"
	"github.com/2451965602/LMS/pkg/errno"
)

func WithLoginData(ctx context.Context, uid int64) context.Context {
	return newContext(ctx, constants.LoginDataKey, strconv.FormatInt(uid, 10))
}

func GetLoginData(ctx context.Context) (int64, error) {
	user, ok := fromContext(ctx, constants.LoginDataKey)
	if !ok {
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get header in context")
	}

	value, err := strconv.ParseInt(user, 10, 64)
	if err != nil {
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to get header in context when parse loginData")
	}
	return value, nil
}
