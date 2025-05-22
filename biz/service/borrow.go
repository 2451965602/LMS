package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/borrow"
	contextLogin "github.com/2451965602/LMS/pkg/base/context"
)

type BorrowService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewBorrowService(ctx context.Context, c *app.RequestContext) *BorrowService {
	return &BorrowService{
		ctx: ctx,
		c:   c,
	}
}

func (s *BorrowService) BookBorrow(ctx context.Context, req borrow.BorrowRequest) (int64, error) {
	userId, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return -1, err
	}

	borrowId, err := db.BookBorrow(ctx, userId, req.BookID)
	if err != nil {
		return -1, err
	}
	return borrowId, nil
}

func (s *BorrowService) BookReturn(ctx context.Context, req borrow.ReturnRequest) (*db.BorrowRecord, error) {
	userId, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, err
	}
	borrowRecord, err := db.BookReturn(ctx, userId, req.BookID, req.BorrowID, req.Status, req.LateFee)
	if err != nil {
		return nil, err
	}
	return borrowRecord, nil
}

func (s *BorrowService) BookRenew(ctx context.Context, req borrow.RenewRequest) (*db.BorrowRecord, error) {
	userId, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, err
	}
	borrowRecord, err := db.BookRenew(ctx, userId, req.BorrowID, int(req.AddTime))
	if err != nil {
		return nil, err
	}
	return borrowRecord, nil
}

func (s *BorrowService) GetCurrentBorrowRecord(ctx context.Context, req borrow.GetBorrowRecordRequest) ([]*db.BorrowRecord, int64, error) {
	userId, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, 0, err
	}

	records, total, err := db.GetCurrentBorrowRecord(ctx, userId, req.PageNum, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var resultRecords []*db.BorrowRecord
	for i := range records {
		resultRecords = append(resultRecords, &records[i])
	}
	return resultRecords, total, nil
}
