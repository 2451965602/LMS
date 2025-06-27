package db

import (
	"context"
	"errors"

	"github.com/2451965602/LMS/biz/model/user"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/pkg/crypt"
	"github.com/2451965602/LMS/pkg/errno"
)

// LoginUser 处理用户登录
// 1. 根据用户名查询用户信息。
// 2. 如果用户不存在，返回错误。
// 3. 验证密码是否正确。
// 4. 如果密码正确，返回用户信息，否则返回错误。
func LoginUser(ctx context.Context, username, password string) (*User, error) {
	var u User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("name = ?", username).
		First(&u).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceUserNotExist, "user not found or invalid credentials (username: %s)", username)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "get user for login failed: %v", err)
	}

	if !crypt.PasswordVerify(password, u.Password) {
		return nil, errno.Errorf(errno.ServiceUserNotExist, "user not found or invalid credentials (password mismatch for username: %s)", username)
	}

	return &u, nil
}

// RegisterUser 注册新用户
// 1. 对密码进行哈希处理。
// 2. 创建新的用户实例并填充请求参数。
// 3. 将用户信息插入到数据库。
// 4. 如果插入成功，返回用户的 ID，否则返回错误。
func RegisterUser(ctx context.Context, username, password, phone string) (int64, error) {
	hashedPassword, err := crypt.PasswordHash(password)
	if err != nil {
		return 0, errno.Errorf(errno.InternalPasswordCryptErrorCode, "encrypt password failed: %v", err)
	}

	u := User{
		Name:       username,
		Password:   hashedPassword,
		Phone:      &phone,
		Permission: "member",
		Status:     "active",
	}

	err = db.WithContext(ctx).
		Table(User{}.TableName()).
		Create(&u).
		Error
	if err != nil {
		return 0, errno.Errorf(errno.InternalDatabaseErrorCode, "create user failed: %v (possible duplicate username '%s')", err, username)
	}
	return u.ID, nil
}

// UpdateUser 更新用户信息
// 1. 根据用户 ID 查询用户信息。
// 2. 如果用户不存在，返回错误。
// 3. 根据请求参数构建更新字段。
// 4. 更新用户信息。
// 5. 如果更新成功，返回更新后的用户信息，否则返回错误。
func UpdateUser(ctx context.Context, userId int64, req user.UpdateUserRequest) (*User, error) {
	var u User
	if errDb := db.WithContext(ctx).Table(User{}.TableName()).Where("id = ?", userId).First(&u).Error; errDb != nil {
		if errors.Is(errDb, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceUserNotExist, "user (id: %d) not exist for update", userId)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch user for update: %v", errDb)
	}

	updates := make(map[string]interface{})
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := crypt.PasswordHash(*req.Password)
		if err != nil {
			return nil, errno.Errorf(errno.InternalPasswordCryptErrorCode, "encrypt password failed for update: %v", err)
		}
		updates["password"] = hashedPassword
		u.Password = hashedPassword
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
		u.Phone = req.Phone
	}

	if len(updates) == 0 {
		return nil, errno.Errorf(errno.ParamMissingErrorCode, "no fields to update")
	}

	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", userId).
		Updates(updates).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update user (id: %d) failed: %v", userId, err)
	}

	return &u, nil
}

// GetUserById 根据用户 ID 获取用户信息
// 1. 根据用户 ID 查询用户信息。
// 2. 如果用户存在，返回用户信息，否则返回错误。
func GetUserById(ctx context.Context, userId int64) (*User, error) {
	var u User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", userId).
		First(&u).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceUserNotExist, "user (id: %d) not exist", userId)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "get user by id (id: %d) failed: %v", userId, err)
	}
	return &u, nil
}

// GetUserByName 根据用户名获取用户信息
// 1. 根据用户名查询用户信息。
// 2. 如果用户存在，返回用户信息，否则返回错误。
func GetUserByName(ctx context.Context, username string) (*User, error) {
	var u User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("name = ?", username).
		First(&u).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceUserNotExist, "user (name: %s) not exist", username)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "get user by name (name: %s) failed: %v", username, err)
	}
	return &u, nil
}

// IsUserExist 检查用户是否存在
// 1. 根据用户名查询用户数量。
// 2. 如果数量大于 0，返回 true，否则返回 false。
func IsUserExist(ctx context.Context, username string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("name = ?", username).
		Count(&count).
		Error
	if err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "check user existence (name: %s) failed: %v", username, err)
	}
	return count > 0, nil
}

// DeleteUser 删除用户
// 1. 根据用户 ID 和用户名验证用户是否存在。
// 2. 如果用户存在，从数据库中删除该用户。
// 3. 如果删除成功，返回 nil，否则返回错误。
func DeleteUser(ctx context.Context, userId int64, username string) error {
	var u User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ? AND name = ?", userId, username).
		First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.Errorf(errno.ServiceUserNotExist, "user with id %d and name '%s' not found for deletion", userId, username)
		}
		return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to verify user for deletion: %v", err)
	}

	result := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ? AND name = ?", userId, username).
		Delete(&User{})

	if result.Error != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "delete user (id: %d, name: %s) failed: %v", userId, username, result.Error)
	}
	if result.RowsAffected == 0 {
		return errno.Errorf(errno.ServiceUserNotExist, "user (id: %d, name: %s) not found during delete operation, or mismatch in ID and name", userId, username)
	}
	return nil
}

// AdminUpdateUser 管理员更新用户信息
// 1. 根据用户 ID 查询用户信息。
// 2. 如果用户不存在，返回错误。
// 3. 根据请求参数构建更新字段。
// 4. 更新用户信息。
// 5. 如果更新成功，返回更新后的用户信息，否则返回错误。
func AdminUpdateUser(ctx context.Context, req user.AdminUpdateUserRequest) (*User, error) {
	var u User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", req.UserID).
		First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceUserNotExist, "user (id: %d) not exist for admin update", req.UserID)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "get user for admin update failed: %v", err)
	}

	updates := make(map[string]interface{})
	if req.Username != nil && *req.Username != "" {
		if *req.Username != u.Name {
			var count int64
			db.WithContext(ctx).Table(User{}.TableName()).Where("name = ? AND id != ?", *req.Username, req.UserID).Count(&count)
			if count > 0 {
				return nil, errno.Errorf(errno.ServiceUserExist, "username '%s' is already taken", *req.Username)
			}
		}
		updates["name"] = *req.Username
		u.Name = *req.Username
	}
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := crypt.PasswordHash(*req.Password)
		if err != nil {
			return nil, errno.Errorf(errno.InternalPasswordCryptErrorCode, "encrypt password for admin update failed: %v", err)
		}
		updates["password"] = hashedPassword
		u.Password = hashedPassword
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
		u.Phone = req.Phone
	}
	if req.Permission != nil {
		updates["permissions"] = *req.Permission
		u.Permission = *req.Permission
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		u.Status = *req.Status
	}

	if len(updates) == 0 {
		return nil, errno.Errorf(errno.ParamMissingErrorCode, "no fields to update for admin")
	}

	err = db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", req.UserID).
		Updates(updates).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "admin update user (id: %d) failed: %v", req.UserID, err)
	}
	return &u, nil
}

// AdminDeleteUser 管理员删除用户
// 1. 根据用户 ID 查询用户信息。
// 2. 如果用户不存在，返回错误。
// 3. 检查用户是否有活动的借阅记录，如果有，不允许删除。
// 4. 如果用户是管理员或图书管理员，不允许删除。
// 5. 如果用户存在且没有活动借阅记录，从数据库中删除该用户。
// 6. 如果删除成功，返回 nil，否则返回错误。
func AdminDeleteUser(ctx context.Context, userId int64) error {
	var userToDelete User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", userId).
		Select("id", "permission").
		First(&userToDelete).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.Errorf(errno.ServiceUserNotExist, "user not exist for admin deletion, ID: %d", userId)
		}
		return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch user for admin deletion: %v", err)
	}

	if userToDelete.Permission == "admin" || userToDelete.Permission == "librarian" {
		return errno.Errorf(errno.ServiceActionNotAllowed, "cannot delete admin or librarian user, ID: %d, Permission: %s", userId, userToDelete.Permission)
	}

	var activeBorrowings int64
	err = db.WithContext(ctx).
		Table(BorrowRecord{}.TableName()).
		Where("user_id = ? AND status IN (?)", userId, []string{"checked_out", "overdue"}).
		Count(&activeBorrowings).Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to check active borrowings for user, ID: %d: %v", userId, err)
	}
	if activeBorrowings > 0 {
		return errno.Errorf(errno.ServiceActionNotAllowed, "user has %d active borrowings (checked_out or overdue), cannot delete, ID: %d", activeBorrowings, userId)
	}

	result := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", userId).
		Delete(&User{})
	if result.Error != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "admin delete user (id: %d) failed: %v", userId, result.Error)
	}
	if result.RowsAffected == 0 {
		return errno.Errorf(errno.ServiceUserNotExist, "user not exist for admin deletion (or already deleted), ID: %d", userId)
	}
	return nil
}

// IsPermission 检查用户是否具有指定权限
// 1. 根据用户 ID 查询用户信息。
// 2. 如果用户不存在，返回错误。
// 3. 根据所需权限检查用户权限。
// 4. 返回检查结果。
func IsPermission(ctx context.Context, userId int64, requiredPermission string) (bool, error) {
	var u User
	err := db.WithContext(ctx).
		Table(User{}.TableName()).
		Where("id = ?", userId).
		Select("permission").
		First(&u).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errno.Errorf(errno.ServiceUserNotExist, "user (id: %d) not found for permission check", userId)
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "check user permission (id: %d) failed: %v", userId, err)
	}

	if requiredPermission == "admin" {
		return u.Permission == "admin", nil
	}
	if requiredPermission == "librarian" {
		return u.Permission == "admin" || u.Permission == "librarian", nil
	}
	if requiredPermission == "member" {
		return u.Permission == "admin" || u.Permission == "librarian" || u.Permission == "member", nil
	}

	return u.Permission == requiredPermission, nil
}
