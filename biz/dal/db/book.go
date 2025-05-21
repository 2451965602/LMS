package db

import (
	"context"
	"time"

	"github.com/2451965602/LMS/biz/service"
	"github.com/2451965602/LMS/pkg/errno"
)

func AddBook(ctx context.Context, book service.Book) (int64, error) {
	bk := Book{
		ISBN:          book.ISBN,
		Location:      book.Location,
		Status:        book.Status,
		PurchasePrice: book.PurchasePrice,
		PurchaseDate:  book.PurchaseDate,
	}
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Create(&bk).
		Error
	if err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "create book failed")
	}
	return bk.ID, nil
}

func UpdateBook(ctx context.Context, bookId int64, location, status *string, purchaseDate *time.Time, purchasePrice *float64) (*Book, error) {
	bk := Book{
		ID: bookId,
	}
	if location != nil {
		bk.Location = *location
	}
	if status != nil {
		bk.Status = *status
	}
	if purchaseDate != nil {
		bk.PurchaseDate = purchaseDate
	}
	if purchasePrice != nil {
		bk.PurchasePrice = purchasePrice
	}
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("ID = ?", bk.ID).
		Updates(&bk).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update book failed")
	}
	return &bk, nil
}

func DeleteBook(ctx context.Context, book service.Book) error {
	err := db.WithContext(ctx).
		Table(Book{}.TableName()).
		Where("ID = ?", book.ID).
		Delete(&Book{}).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "delete book failed")
	}
	return nil
}

func SearchBook(ctx context.Context, isbn *string, bookId *int64) ([]Book, error) {
	var results []Book
	query := db.WithContext(ctx).Table(Book{}.TableName())
	if isbn != nil {
		query = query.Where("ISBN = ?", *isbn)
	}
	if bookId != nil {
		query = query.Where("ID = ?", *bookId)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "search book failed")
	}

	if len(results) == 0 {
		return nil, errno.NewErrNo(errno.ServiceBookNotExist, "book not exist")
	}

	return results, nil
}
