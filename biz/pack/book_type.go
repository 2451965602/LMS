package pack

import (
	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/model"
)

func BuildBookTypeResp(info *db.BookType) *model.BookType {
	if info == nil {
		return nil
	}
	return &model.BookType{
		ISBN:            info.ISBN,
		Title:           info.Title,
		Author:          info.Author,
		Category:        info.Category,
		Publisher:       info.Publisher,
		PublishYear:     info.PublishYear,
		Description:     info.Description,
		TotalCopies:     info.TotalCopies,
		AvailableCopies: info.AvailableCopies,
	}
}

func BuildBookTypeListResp(infos []*db.BookType) []*model.BookType {
	if infos == nil {
		return nil
	}
	resp := make([]*model.BookType, 0, len(infos))
	for _, info := range infos {
		resp = append(resp, BuildBookTypeResp(info))
	}
	return resp
}
