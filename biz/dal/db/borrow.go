package db

import (
	"context"
	"errors"
	"time"

	"github.com/2451965602/LMS/pkg/constants"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/pkg/errno"
)

func BookBorrow(ctx context.Context, userId, bookId int64) (int64, error) {
	var br BorrowRecord
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var bookInfo Book
		if err := tx.Table(Book{}.TableName()).Where("id = ?", bookId).First(&bookInfo).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.NewErrNo(errno.ServiceBookNotExist, "book not exist")
			}
			return errno.Errorf(errno.InternalDatabaseErrorCode, "get book info failed for borrow: %v", err)
		}

		if bookInfo.Status != "available" {
			return errno.Errorf(errno.ServiceActionNotAllowed, "book with id %d is not available (status: %s)", bookId, bookInfo.Status)
		}

		var bt BookType
		if err := tx.Table(BookType{}.TableName()).Where("ISBN = ?", bookInfo.ISBN).First(&bt).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch book type %s: %v", bookInfo.ISBN, err)
		}
		if bt.AvailableCopies <= 0 {
			return errno.Errorf(errno.ServiceActionNotAllowed, "no available copies for book type %s (ISBN: %s)", bt.Title, bt.ISBN)
		}

		br = BorrowRecord{
			UserID:       userId,
			BookID:       bookId,
			CheckoutDate: time.Now(),
			DueDate:      time.Now().AddDate(0, 0, constants.DefauteRenewTime),
			Status:       "checked_out",
			RenewalCount: 0,
		}
		if err := tx.Table(BorrowRecord{}.TableName()).Create(&br).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "create borrow record failed: %v", err)
		}

		result := tx.Table(BookType{}.TableName()).
			Where("ISBN = ?", bookInfo.ISBN).
			Update("available_copies", gorm.Expr("available_copies - 1"))
		if result.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "sub book available count failed: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to update available copies for book type %s (ISBN not found or no change)", bookInfo.ISBN)
		}

		if err := tx.Table(Book{}.TableName()).Where("id = ?", bookId).Update("status", "checked_out").Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "update book status to checked_out failed: %v", err)
		}

		return nil
	})
	if err != nil {
		return -1, err
	}
	return br.ID, nil
}

func BookReturn(ctx context.Context, userId, bookId, borrowId int64, returnStatus string, lateFee float64) (*BorrowRecord, error) {
	var updatedBr BorrowRecord

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var bookInfo Book
		if err := tx.Table(Book{}.TableName()).Where("id = ?", bookId).First(&bookInfo).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.NewErrNo(errno.ServiceBookNotExist, "book not exist")
			}
			return errno.Errorf(errno.InternalDatabaseErrorCode, "get book info failed for return: %v", err)
		}

		var currentBr BorrowRecord
		err := tx.Table(BorrowRecord{}.TableName()).
			Where("id = ? AND user_id = ? AND book_id = ?", borrowId, userId, bookId).
			First(&currentBr).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.Errorf(errno.ServiceBorrowRecordNotExist,
					"borrow record not found or does not match user/book (id: %d, user_id: %d, book_id: %d)",
					borrowId, userId, bookId)
			}
			return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch borrow record: %v", err)
		}

		if currentBr.Status == "returned" || currentBr.Status == "lost" {
			return errno.Errorf(errno.ServiceActionNotAllowed, "book already in terminal status: %s", currentBr.Status)
		}

		currentTime := time.Now()
		updates := map[string]interface{}{
			"status":      returnStatus,
			"late_fee":    lateFee,
			"return_date": &currentTime,
		}
		result := tx.Table(BorrowRecord{}.TableName()).
			Where("id = ? AND user_id = ?", borrowId, userId).
			Updates(updates)

		if result.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "update borrow record failed: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			return errno.Errorf(errno.ServiceBorrowRecordNotExist, "borrow record not found or not updated during return process")
		}

		if returnStatus == "returned" {
			resultInc := tx.Table(BookType{}.TableName()).
				Where("ISBN = ?", bookInfo.ISBN).
				Update("available_copies", gorm.Expr("available_copies + 1"))
			if resultInc.Error != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "add book available count failed: %v", resultInc.Error)
			}
			if err := tx.Table(Book{}.TableName()).Where("id = ?", bookId).Update("status", "available").Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "update book status to available failed: %v", err)
			}
		} else if returnStatus == "lost" || returnStatus == "damaged" {
			if err := tx.Table(Book{}.TableName()).Where("id = ?", bookId).Update("status", returnStatus).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "update book status to %s failed: %v", returnStatus, err)
			}
		}

		if err := tx.Table(BorrowRecord{}.TableName()).Where("id = ?", borrowId).First(&updatedBr).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch updated borrow record: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updatedBr, nil
}

func BookRenew(ctx context.Context, userId, borrowId int64, daysToExtend int) (*BorrowRecord, error) {
	var record BorrowRecord

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(BorrowRecord{}.TableName()).
			Where("id = ? AND user_id = ?", borrowId, userId).
			First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.Errorf(errno.ServiceBorrowRecordNotExist, "borrow record not found (id: %d) or does not belong to user (user_id: %d)", borrowId, userId)
			}
			return errno.Errorf(errno.InternalDatabaseErrorCode, "get borrow record failed for renew: %v", err)
		}

		if record.Status != "checked_out" {
			return errno.Errorf(errno.ServiceActionNotAllowed, "cannot renew book, status is '%s', not 'checked_out'", record.Status)
		}
		if record.RenewalCount >= constants.MaxRenewTime {
			return errno.Errorf(errno.ServiceActionNotAllowed, "cannot renew book, maximum renewal count (2) reached (current: %d)", record.RenewalCount)
		}

		newDueDate := record.DueDate.AddDate(0, 0, daysToExtend)

		updateResult := tx.Table(BorrowRecord{}.TableName()).
			Where("id = ?", borrowId).
			Updates(map[string]interface{}{
				"due_date":      newDueDate,
				"renewal_count": gorm.Expr("renewal_count + 1"),
			})
		if updateResult.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "update borrow record for renew failed: %v", updateResult.Error)
		}
		if updateResult.RowsAffected == 0 {
			return errno.Errorf(errno.ServiceBorrowRecordNotExist, "borrow record (id: %d) not found during renewal update", borrowId)
		}

		if err := tx.Table(BorrowRecord{}.TableName()).
			Where("id = ?", borrowId).First(&record).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch updated record post-renewal: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func GetCurrentBorrowRecord(ctx context.Context, userId, pageNum, pageSize, status int64) ([]BorrowRecord, int64, error) {
	var results []BorrowRecord
	var total int64

	baseQuery := db.WithContext(ctx).Table(BorrowRecord{}.TableName()).Where("user_id = ?", userId)

	switch status {
	case constants.CheckedOut:
		baseQuery = baseQuery.Where("status = ?", "checked_out")
	case constants.Returned:
		baseQuery = baseQuery.Where("status = ?", "returned")
	case constants.Overdue:
		baseQuery = baseQuery.Where("status = ?", "overdue")
	case constants.Lost:
		baseQuery = baseQuery.Where("status = ?", "lost")
	}

	err := baseQuery.Count(&total).Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "count borrow records failed: %v", err)
	}

	if total == 0 {
		return []BorrowRecord{}, 0, nil
	}

	offset := int((pageNum - 1) * pageSize)
	if offset < 0 {
		offset = 0
	}

	err = baseQuery.Order("checkout_date DESC").
		Offset(offset).
		Limit(int(pageSize)).
		Find(&results).
		Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "search borrow record failed: %v", err)
	}
	return results, total, nil
}

func GetBorrowRecordById(ctx context.Context, borrowId int64) (*BorrowRecord, error) {
	var record BorrowRecord
	err := db.WithContext(ctx).
		Table(BorrowRecord{}.TableName()).
		Where("id = ?", borrowId).
		First(&record).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceBorrowRecordNotExist, "borrow record not found with ID: %d", borrowId)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "get borrow record by id failed: %v", err)
	}
	return &record, nil
}
