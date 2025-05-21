package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/biz/service"
	"github.com/2451965602/LMS/pkg/errno"
)

func AddBookType(ctx context.Context, bookType service.BookType) (*BookType, error) {
	bt := BookType{
		Title:       bookType.Title,
		Author:      bookType.Author,
		Category:    bookType.Category,
		ISBN:        bookType.ISBN,
		Publisher:   bookType.Publisher,
		PublishYear: bookType.PublishYear,
		Description: bookType.Description,
	}
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Create(&bt).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "create book type failed")
	}
	return &bt, nil
}

func UpdateBookType(ctx context.Context, isbn string, title, author, category, publisher, description *string, publishYear *int) (BookType, error) {
	bt := BookType{
		ISBN: isbn,
	}
	if title != nil {
		bt.Title = *title
	}
	if author != nil {
		bt.Author = *author
	}
	if category != nil {
		bt.Category = *category
	}
	if publisher != nil {
		bt.Publisher = *publisher
	}
	if publishYear != nil {
		bt.PublishYear = publishYear
	}
	if description != nil {
		bt.Description = *description
	}
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", bt.ISBN).
		Updates(&bt).
		Error
	if err != nil {
		return BookType{}, errno.Errorf(errno.InternalDatabaseErrorCode, "update book type failed")
	}
	return bt, nil
}

func DeleteBookType(ctx context.Context, isbn string) error {
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", isbn).
		Delete(&BookType{}).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "delete book type failed")
	}
	return nil
}

func SearchBookType(ctx context.Context, title, isbn, author *string) ([]BookType, error) {
	var results []BookType

	query := db.WithContext(ctx).Table(BookType{}.TableName())
	if title != nil {
		query = query.Where("Title = ?", *title)
	}
	if isbn != nil {
		query = query.Where("ISBN = ?", *isbn)
	}
	if author != nil {
		query = query.Where("Author = ?", *author)
	}
	err := query.Find(&results).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "search book type failed")
	}

	if len(results) == 0 {
		return nil, errno.Errorf(errno.ServiceBookTypeNotFound, "book type not found")
	}
	return results, nil
}

func AddBookTotalCount(ctx context.Context, isbn string) error {
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", isbn).
		Update("TotalCopies", gorm.Expr("TotalCopies + 1")).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "add book count failed")
	}
	return nil
}

func SubBookTotalCount(ctx context.Context, ISBN string) error {
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", ISBN).
		Update("TotalCopies", gorm.Expr("TotalCopies - 1")).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "sub book count failed")
	}
	return nil
}

func addBookAvailableCount(ctx context.Context, ISBN string) error {
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", ISBN).
		Update("AvailableCopies", gorm.Expr("AvailableCopies + 1")).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "add book available count failed")
	}
	return nil
}

func subBookAvailableCount(ctx context.Context, ISBN string) error {
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", ISBN).
		Update("AvailableCopies", gorm.Expr("AvailableCopies - 1")).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "sub book available count failed")
	}
	return nil
}
