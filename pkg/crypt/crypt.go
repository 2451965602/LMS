package crypt

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/2451965602/LMS/pkg/errno"
)

// PasswordHash 对密码进行哈希加密
// 参数：
//   - pwd: 明文密码
//
// 返回值：
//   - string: 加密后的密码哈希值
//   - error: 错误信息，如果加密失败会返回错误
func PasswordHash(pwd string) (string, error) {
	// 使用bcrypt算法对密码进行哈希加密，使用默认的加密成本
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		// 如果加密失败，返回错误
		return "", errno.NewErrNo(errno.InternalPasswordCryptErrorCode, "encrypt password failed")
	}

	// 返回加密后的密码哈希值
	return string(bytes), nil
}

// PasswordVerify 验证密码是否与哈希值匹配
// 参数：
//   - pwd: 明文密码
//   - hash: 加密后的密码哈希值
//
// 返回值：
//   - bool: 如果密码与哈希值匹配返回true，否则返回false
func PasswordVerify(pwd, hash string) bool {
	// 使用bcrypt算法比较明文密码和哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))

	// 如果没有错误，说明密码匹配
	return err == nil
}
