package pack

import (
	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/model"
)

func BuildBookResp(info *db.Book) *model.Book {
	if info == nil {
		return nil
	}
	result := &model.Book{
		ID:            info.ID,
		Isbn:          info.ISBN,
		Location:      info.Location,
		Status:        info.Status,
		PurchasePrice: info.PurchasePrice,
		PurchaseDate:  info.PurchaseDate.Format("2006-01-02 15:04:05"),
	}
	if info.LastCheckout != nil {
		result.LastCheckout = info.LastCheckout.Format("2006-01-02 15:04:05")
	} else {
		result.LastCheckout = "" // or some default value
	}

	return result
}

func BuildBookListResp(infos []*db.Book) []*model.Book {
	if infos == nil {
		return nil
	}
	resp := make([]*model.Book, 0, len(infos))
	for _, info := range infos {
		resp = append(resp, BuildBookResp(info))
	}
	return resp
}
