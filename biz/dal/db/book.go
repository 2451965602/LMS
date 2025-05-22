package db

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/biz/model/book"
	"github.com/2451965602/LMS/pkg/errno"
)

func AddBook(ctx context.Context, req book.AddBookRequest) (int64, error) {
	bk := Book{
		ISBN:          req.ISBN,
		Location:      req.Location,
		Status:        req.Status,
		PurchasePrice: req.PurchasePrice,
		PurchaseDate:  time.Unix(req.PurchaseDate, 0),
	}

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(Book{}.TableName()).Create(&bk).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "create book failed: %v", err)
		}

		result := tx.Table(BookType{}.TableName()).
			Where("ISBN = ?", bk.ISBN).
			Updates(map[string]interface{}{
				"total_copies":     gorm.Expr("total_copies + 1"),
				"available_copies": gorm.Expr("available_copies + 1"),
			})

		if result.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "add book count failed: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			return errno.Errorf(errno.ServiceBookTypeNotFound, "book type with ISBN %s not found, cannot update copies", bk.ISBN)
		}
		return nil
	})
	if err != nil {
		return -1, err
	}

	return bk.ID, nil
}

func UpdateBook(ctx context.Context, req book.UpdateBookRequest) (*Book, error) {
	var bk Book
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("id = ?", req.BookID).
		First(&bk).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceBookNotExist, "book not exist")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch book for update: %v", err)
	}

	updates := make(map[string]interface{})
	if req.Location != nil {
		updates["location"] = *req.Location
		bk.Location = *req.Location
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		bk.Status = *req.Status
	}
	if req.PurchaseDate != nil {
		updates["purchase_date"] = time.Unix(*req.PurchaseDate, 0)
		bk.PurchaseDate = time.Unix(*req.PurchaseDate, 0)
	}
	if req.PurchasePrice != nil {
		updates["purchase_price"] = *req.PurchasePrice
		bk.PurchasePrice = *req.PurchasePrice
	}

	if len(updates) == 0 {
		return nil, errno.Errorf(errno.ParamMissingErrorCode, "no fields to update")
	}

	err = db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("id = ?", req.BookID).
		Updates(updates).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update book failed: %v", err)
	}
	return &bk, nil
}

func DeleteBook(ctx context.Context, bookId int64) error {
	var bk Book
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("id = ?", bookId).
		First(&bk).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NewErrNo(errno.ServiceBookNotExist, "book not exist, cannot delete")
		}
		return errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch book for deletion: %v", err)
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(Book{}.TableName()).
			Where("id = ?", bookId).
			Delete(&Book{})

		if result.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "delete book failed: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			return errno.NewErrNo(errno.ServiceBookNotExist, "book not found during delete operation")
		}

		availableExpr := "available_copies - 1"
		if bk.Status != "available" {
			availableExpr = "available_copies"
		}

		updateResult := tx.Table(BookType{}.TableName()).
			Where("ISBN = ?", bk.ISBN).
			Updates(map[string]interface{}{
				"total_copies":     gorm.Expr("total_copies - 1"),
				"available_copies": gorm.Expr(availableExpr),
			})

		if updateResult.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "sub book count failed: %v", updateResult.Error)
		}
		return nil
	})

	return err
}

func SearchBook(ctx context.Context, req book.GetBookRequest) ([]*Book, int64, error) {
	var results []Book
	var total int64

	query := db.WithContext(ctx).Table(Book{}.TableName())
	countQuery := db.WithContext(ctx).Table(Book{}.TableName())

	if req.ISBN != nil && *req.ISBN != "" {
		query = query.Where("ISBN = ?", *req.ISBN)
		countQuery = countQuery.Where("ISBN = ?", *req.ISBN)
	}
	if req.BookID != nil {
		query = query.Where("id = ?", *req.BookID)
		countQuery = countQuery.Where("id = ?", *req.BookID)
	}

	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "count books failed: %v", err)
	}

	if total == 0 {
		return []*Book{}, 0, nil
	}

	offset := int((req.PageNum - 1) * req.PageSize)
	if offset < 0 {
		offset = 0
	}

	err = query.Order("id desc").
		Offset(offset).
		Limit(int(req.PageSize)).
		Find(&results).
		Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "search book failed: %v", err)
	}

	var resultBooks []*Book
	for i := range results {
		resultBooks = append(resultBooks, &results[i])
	}

	return resultBooks, total, nil
}

func IsBookExist(ctx context.Context, bookId int64) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("id = ?", bookId).
		Count(&count).
		Error
	if err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "check book existence failed: %v", err)
	}
	return count > 0, nil
}

func GetBookById(ctx context.Context, bookId int64) (*Book, error) {
	var info Book
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("id = ?", bookId).
		First(&info).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceBookNotExist, "book not exist")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "get book by id failed: %v", err)
	}
	return &info, nil
}

func IsBookInISBN(ctx context.Context, isbn string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("ISBN = ?", isbn).
		Count(&count).
		Error
	if err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "check book existence failed: %v", err)
	}
	return count > 0, nil
}
