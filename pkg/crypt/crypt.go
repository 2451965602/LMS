package crypt

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/2451965602/LMS/pkg/errno"
)

func PasswordHash(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", errno.NewErrNo(errno.InternalPasswordCryptErrorCode, "encrypt password failed")
	}

	return string(bytes), nil
}

func PasswordVerify(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))

	return err == nil
}
