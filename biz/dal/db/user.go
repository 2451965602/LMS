package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/biz/service"
	"github.com/2451965602/LMS/pkg/crypt"
	"github.com/2451965602/LMS/pkg/errno"
)

func LoginUser(ctx context.Context, username, password string) (*User, error) {
	user, err := GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}

	if !crypt.PasswordVerify(password, user.Password) {
		return nil, errno.Errorf(errno.ServiceUserNotExist, "invalid credentials")
	}

	return &user, nil
}

func RegisterUser(ctx context.Context, username, password string) (int64, error) {
	var err error
	user := User{
		Name:     username,
		Password: password,
	}

	user.Password, err = crypt.PasswordHash(user.Password)
	if err != nil {
		return 0, errno.Errorf(errno.InternalPasswordCryptErrorCode, "encrypt password failed")
	}

	err = db.WithContext(ctx).
		Table(User{}.TableName()).
		Create(&user).
		Error
	if err != nil {
		return 0, errno.Errorf(errno.InternalDatabaseErrorCode, "create user failed")
	}
	return user.ID, nil
}

func UpdateUser(ctx context.Context, user service.User) (*User, error) {
	u := User{}
	if user.Name != "" {
		u.Name = user.Name
	}
	if user.Password != "" {
		var err error
		u.Password, err = crypt.PasswordHash(user.Password)
		if err != nil {
			return nil, errno.Errorf(errno.InternalPasswordCryptErrorCode, "encrypt password failed")
		}
	}
	if user.Phone != nil {
		u.Phone = user.Phone
	}

	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", u.ID).
		Updates(&u).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update user failed")
	}
	return &u, nil
}

func GetUserById(ctx context.Context, userId int64) (User, error) {
	var user User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", userId).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, errno.Errorf(errno.ServiceUserNotExist, "user not exist")
		}
		return User{}, errno.Errorf(errno.InternalDatabaseErrorCode, "get user failed")
	}
	return user, nil
}

func GetUserByName(ctx context.Context, username string) (User, error) {
	var user User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("name = ?", username).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, errno.Errorf(errno.ServiceUserNotExist, "user not exist")
		}
		return User{}, errno.Errorf(errno.InternalDatabaseErrorCode, "get user failed")
	}
	return user, nil
}

func IsUserExist(ctx context.Context, username string) (bool, error) {
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("username = ?", username).
		First(&User{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "check user exist failed")
	}
	return true, nil
}

func DeleteUser(ctx context.Context, userId int64, username string) error {
	user := User{
		ID:   userId,
		Name: username,
	}
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ? AND name = ?", user.ID, user.Name).
		Delete(&user).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "delete user failed")
	}
	return nil
}

func AdminUpdateUser(ctx context.Context, user service.User) (*User, error) {
	u := User{}
	if user.Name != "" {
		u.Name = user.Name
	}
	if user.Password != "" {
		var err error
		u.Password, err = crypt.PasswordHash(user.Password)
		if err != nil {
			return nil, errno.Errorf(errno.InternalPasswordCryptErrorCode, "encrypt password failed")
		}
	}
	if user.Phone != nil {
		u.Phone = user.Phone
	}
	if user.Permissions != "" {
		u.Permissions = user.Permissions
	}
	if user.Status != "" {
		u.Status = user.Status
	}

	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", u.ID).
		Updates(&u).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update user failed")
	}
	return &u, nil
}

func AdminDeleteUser(ctx context.Context, userId string) error {
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", userId).
		Delete(&User{}).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "delete user failed")
	}
	return nil
}

func IsPermission(ctx context.Context, userId int64, op string) (bool, error) {
	var user User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ? ", userId).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "check user permission failed")
	}

	if user.Permissions != op {
		return false, nil
	}
	return true, nil
}
