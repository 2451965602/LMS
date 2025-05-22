package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/biz/model/booktype"
	"github.com/2451965602/LMS/pkg/errno"
)

func AddBookType(ctx context.Context, req booktype.AddBookTypeRequest) (*BookType, error) {
	bt := BookType{
		Title:       req.Title,
		Author:      req.Author,
		Category:    req.Category,
		ISBN:        req.ISBN,
		Publisher:   req.Publisher,
		PublishYear: req.PublishYear,
		Description: req.Description,
	}
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Create(&bt).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "create book type failed: %v", err)
	}
	return &bt, nil
}

func UpdateBookType(ctx context.Context, req booktype.UpdateBookTypeRequest) (*BookType, error) {
	var bt BookType

	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", req.ISBN).
		First(&bt).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceBookTypeNotFound, "book type with ISBN %s not found for update", req.ISBN)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch book type for update: %v", err)
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
		bt.Title = *req.Title
	}
	if req.Author != nil {
		updates["author"] = *req.Author
		bt.Author = *req.Author
	}
	if req.Category != nil {
		updates["category"] = *req.Category
		bt.Category = *req.Category
	}
	if req.Publisher != nil {
		updates["publisher"] = *req.Publisher
		bt.Publisher = *req.Publisher
	}
	if req.PublishYear != nil {
		updates["publish_year"] = *req.PublishYear
		bt.PublishYear = *req.PublishYear
	}
	if req.Description != nil {
		updates["description"] = *req.Description
		bt.Description = *req.Description
	}

	if len(updates) == 0 {
		return nil, errno.Errorf(errno.ParamMissingErrorCode, "no fields to update")
	}

	err = db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", req.ISBN).
		Updates(updates).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "update book type failed: %v", err)
	}

	return &bt, nil
}

func DeleteBookType(ctx context.Context, isbn string) error {
	result := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", isbn).
		Delete(&BookType{})

	if result.Error != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "delete book type failed: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errno.Errorf(errno.ServiceBookTypeNotFound, "book type not found with ISBN: %s, no rows deleted", isbn)
	}
	return nil
}

func SearchBookType(ctx context.Context, title, author *string, isbn *string, pageNum, pageSize int64) ([]*BookType, int64, error) {
	var results []BookType
	var total int64

	query := db.WithContext(ctx).Table(BookType{}.TableName())
	countQuery := db.WithContext(ctx).Table(BookType{}.TableName())

	if title != nil && *title != "" {
		query = query.Where("title LIKE ?", "%"+*title+"%")
		countQuery = countQuery.Where("title LIKE ?", "%"+*title+"%")
	}
	if isbn != nil && *isbn != "" {
		query = query.Where("ISBN = ?", *isbn)
		countQuery = countQuery.Where("ISBN = ?", *isbn)
	}
	if author != nil && *author != "" {
		query = query.Where("author LIKE ?", "%"+*author+"%")
		countQuery = countQuery.Where("author LIKE ?", "%"+*author+"%")
	}

	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "count book type failed: %v", err)
	}

	if total == 0 {
		return []*BookType{}, 0, nil
	}

	offset := int((pageNum - 1) * pageSize)
	if offset < 0 {
		offset = 0
	}

	err = query.Order("ISBN desc").
		Offset(offset).
		Limit(int(pageSize)).
		Find(&results).
		Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "search book type failed: %v", err)
	}

	var resultBookTypes []*BookType
	for i := range results {
		resultBookTypes = append(resultBookTypes, &results[i])
	}

	return resultBookTypes, total, nil
}

func IsBookTypeExist(ctx context.Context, isbn string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", isbn).
		Count(&count).
		Error
	if err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "check book type existence failed: %v", err)
	}
	return count > 0, nil
}

func GetBookTypeByISBN(ctx context.Context, isbn string) (*BookType, error) {
	var bt BookType
	err := db.WithContext(ctx).
		Table(BookType{}.TableName()).
		Where("ISBN = ?", isbn).
		First(&bt).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceBookTypeNotFound, "book type with ISBN %s not found", isbn)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "failed to fetch book type: %v", err)
	}
	return &bt, nil
}
