package pack

import (
	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/model"
)

func BuildUserResp(info *db.User) *model.User {
	if info == nil {
		return nil
	}
	return &model.User{
		ID:           info.ID,
		Username:     info.Name,
		Phone:        info.Phone,
		Status:       info.Status,
		Permissions:  info.Permission,
		RegisterDate: info.RegisterDate.Format("2006-01-02"),
	}
}
