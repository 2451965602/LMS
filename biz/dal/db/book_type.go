package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/biz/model/booktype"
	"github.com/2451965602/LMS/pkg/errno"
)

// AddBookType 添加一个新的书籍类型到数据库
// 1. 根据请求参数创建一个新的 BookType 实例。
// 2. 使用 gorm 的 Create 方法插入到数据库。
// 3. 如果插入成功，返回新创建的 BookType 实例，否则返回错误。
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

// UpdateBookType 更新指定 ISBN 的书籍类型信息
// 1. 根据 ISBN 查询书籍类型是否存在。
// 2. 根据请求参数构建更新字段。
// 3. 使用 gorm 的 Updates 方法更新书籍类型信息。
// 4. 如果更新成功，返回更新后的书籍类型信息，否则返回错误。
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

// DeleteBookType 删除指定 ISBN 的书籍类型
// 1. 根据 ISBN 查询书籍类型是否存在。
// 2. 如果存在，从数据库中删除该书籍类型。
// 3. 如果删除成功，返回 nil，否则返回错误。
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

// SearchBookType 搜索书籍类型
// 1. 根据请求参数构建查询条件。
// 2. 查询总记录数。
// 3. 根据分页参数查询书籍类型列表。
// 4. 返回书籍类型列表和总记录数。
func SearchBookType(ctx context.Context, title, author, isbn, category *string, pageNum, pageSize int64) ([]*BookType, int64, error) {
	var results []BookType
	var total int64

	query := db.WithContext(ctx).Table(BookType{}.TableName())
	countQuery := db.WithContext(ctx).Table(BookType{}.TableName())

	// 构建查询条件
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
	if category != nil && *category != "" {
		query = query.Where("category LIKE ?", "%"+*category+"%")
		countQuery = countQuery.Where("category LIKE ?", "%"+*category+"%")
	}

	// 查询总记录数
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "count book type failed: %v", err)
	}

	// 如果没有记录，直接返回空列表
	if total == 0 {
		return []*BookType{}, 0, nil
	}

	// 计算分页的偏移量
	offset := int((pageNum - 1) * pageSize)
	if offset < 0 {
		offset = 0
	}

	// 查询书籍类型列表
	err = query.Order("ISBN desc").
		Offset(offset).
		Limit(int(pageSize)).
		Find(&results).
		Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "search book type failed: %v", err)
	}

	// 将结果转换为指针切片
	var resultBookTypes []*BookType
	for i := range results {
		resultBookTypes = append(resultBookTypes, &results[i])
	}

	return resultBookTypes, total, nil
}

// IsBookTypeExist 检查指定 ISBN 的书籍类型是否存在
// 1. 根据 ISBN 查询书籍类型的数量。
// 2. 如果数量大于 0，返回 true，否则返回 false。
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

// GetBookTypeByISBN 根据 ISBN 获取书籍类型信息
// 1. 根据 ISBN 查询书籍类型。
// 2. 如果书籍类型存在，返回书籍类型信息，否则返回错误。
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
