package pack

import (
	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/model"
)

func BuildBorrowRecordResp(info *db.BorrowRecord) *model.BorrowRecord {
	if info == nil {
		return nil
	}
	result := &model.BorrowRecord{
		ID:           info.ID,
		UserID:       info.UserID,
		BookID:       info.BookID,
		CheckoutDate: info.CheckoutDate.Format("2006-01-02 15:04:05"),
		DueDate:      info.DueDate.Format("2006-01-02 15:04:05"),
		Status:       info.Status,
		RenewalCount: info.RenewalCount,
		LateFee:      info.LateFee,
	}
	if info.ReturnDate != nil {
		result.ReturnDate = info.ReturnDate.Format("2006-01-02 15:04:05")
	} else {
		result.ReturnDate = ""
	}
	return result
}

func BuildBorrowRecordListResp(infos []*db.BorrowRecord) []*model.BorrowRecord {
	if infos == nil {
		return nil
	}
	resp := make([]*model.BorrowRecord, 0, len(infos))
	for _, info := range infos {
		resp = append(resp, BuildBorrowRecordResp(info))
	}
	return resp
}
