package db

import (
	"context"
	"time"

	"github.com/2451965602/LMS/pkg/errno"
)

func BookBorrow(ctx context.Context, userId, bookId int64) (int64, error) {
	br := BorrowRecord{
		UserID: userId,
		BookID: bookId,
	}
	err := db.WithContext(ctx).
		Table(BorrowRecord{}.TableName()).
		Create(&br).
		Error
	if err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "create borrow record failed")
	}
	return br.ID, nil
}

func BookReturn(ctx context.Context, bookId, borrowId int64, statue string, lateFee float64) (*BorrowRecord, error) {
	br := BorrowRecord{
		ID:      borrowId,
		BookID:  bookId,
		Status:  statue,
		LateFee: lateFee,
	}
	err := db.WithContext(ctx).
		Table(BorrowRecord{}.TableName()).
		Where("ID = ? AND UserID = ?", br.ID, br.UserID).
		Updates(&br).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update borrow record failed")
	}
	return &br, nil
}

func BookRenew(ctx context.Context, userId, borrowId int64, newDueDate time.Time) (*BorrowRecord, error) {
	err := db.WithContext(ctx).
		Table(BorrowRecord{}.TableName()).
		Where("ID = ? AND UserID = ?", borrowId, userId).
		Updates(map[string]interface{}{
			"DueDate": newDueDate,
			"Status":  "checked_out",
		}).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update borrow record failed")
	}

	var record BorrowRecord
	err = db.WithContext(ctx).
		Table(BorrowRecord{}.TableName()).
		Where("ID = ?", borrowId).
		First(&record).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "get borrow record failed")
	}

	return &record, nil
}

func GetCurrentBorrowRecord(ctx context.Context, bookId, userId *int64) ([]BorrowRecord, error) {
	var results []BorrowRecord
	query := db.WithContext(ctx).Table(BorrowRecord{}.TableName())
	if bookId != nil {
		query = query.Where("BookID = ?", *bookId)
	}
	if userId != nil {
		query = query.Where("UserID = ?", *userId)
	}
	err := query.Find(&results).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "search borrow record failed")
	}
	if len(results) == 0 {
		return nil, errno.Errorf(errno.ServiceBorrowRecordNotExist, "no borrow record found")
	}
	return results, nil
}

func BookReserve(ctx context.Context, bookId int64, reserveDate *time.Time) (int64, error) {
	br := Reservation{
		BookID:      bookId,
		ReserveDate: *reserveDate,
	}
	err := db.WithContext(ctx).
		Table(Reservation{}.TableName()).
		Create(&br).
		Error
	if err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "create reservation record failed")
	}
	return br.ID, nil
}

func CancelBookReserve(ctx context.Context, userId, reserveId int64) (int64, error) {
	err := db.WithContext(ctx).
		Table(Reservation{}.TableName()).
		Where("ID = ? AND UserID = ?", reserveId, userId).
		Updates(map[string]interface{}{
			"Status": "canceled",
		}).
		Error
	if err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "cancel reservation record failed")
	}

	return reserveId, nil
}

func GetCurrentReservation(ctx context.Context, userId int64) ([]Reservation, error) {
	var results []Reservation
	err := db.WithContext(ctx).
		Table(Reservation{}.TableName()).
		Where("UserID = ?", userId).
		Find(&results).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "search reservation record failed")
	}
	if len(results) == 0 {
		return nil, errno.Errorf(errno.ServiceReservationRecordNotExist, "no reservation record found")
	}
	return results, nil
}
