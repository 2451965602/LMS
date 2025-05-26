package db

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/2451965602/LMS/biz/model/book"
	"github.com/2451965602/LMS/pkg/errno"
)

// AddBook 添加一本新书到数据库中
//  1. 创建一个新的 Book 实例并填充请求参数。
//  2. 使用事务确保操作的原子性：
//     a. 将新书插入到 Book 表中。
//     b. 更新 BookType 表中的总副本数和可用副本数。
//  3. 如果事务成功，返回新书的 ID，否则返回错误。
func AddBook(ctx context.Context, req book.AddBookRequest) (int64, error) {
	bk := Book{
		ISBN:          req.ISBN,
		Location:      req.Location,
		Status:        req.Status,
		PurchasePrice: req.PurchasePrice,
		PurchaseDate:  time.Unix(req.PurchaseDate, 0),
	}

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 插入新书到 Book 表
		if err := tx.Table(Book{}.TableName()).Create(&bk).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "create book failed: %v", err)
		}

		// 更新 BookType 表中的总副本数和可用副本数
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

// UpdateBook 更新指定 ID 的书籍信息
// 1. 根据 BookID 查询书籍是否存在。
// 2. 根据请求参数构建更新字段。
// 3. 使用 gorm 的 Updates 方法更新书籍信息。
// 4. 如果更新成功，返回更新后的书籍信息，否则返回错误。
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

// DeleteBook 删除指定 ID 的书籍
//  1. 根据 BookID 查询书籍是否存在。
//  2. 使用事务确保操作的原子性：
//     a. 从 Book 表中删除书籍。
//     b. 更新 BookType 表中的总副本数和可用副本数。
//  3. 如果事务成功，返回 nil，否则返回错误。
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
		// 从 Book 表中删除书籍
		result := tx.Table(Book{}.TableName()).
			Where("id = ?", bookId).
			Delete(&Book{})

		if result.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "delete book failed: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			return errno.NewErrNo(errno.ServiceBookNotExist, "book not found during delete operation")
		}

		// 更新 BookType 表中的总副本数和可用副本数
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

// SearchBook 搜索书籍
// 1. 根据请求参数构建查询条件。
// 2. 查询总记录数。
// 3. 根据分页参数查询书籍列表。
// 4. 返回书籍列表和总记录数。
func SearchBook(ctx context.Context, req book.GetBookRequest) ([]*Book, int64, error) {
	var results []Book
	var total int64

	query := db.WithContext(ctx).Table(Book{}.TableName())
	countQuery := db.WithContext(ctx).Table(Book{}.TableName())

	// 构建查询条件
	if req.ISBN != nil && *req.ISBN != "" {
		query = query.Where("ISBN = ?", *req.ISBN)
		countQuery = countQuery.Where("ISBN = ?", *req.ISBN)
	}
	if req.BookID != nil {
		query = query.Where("id = ?", *req.BookID)
		countQuery = countQuery.Where("id = ?", *req.BookID)
	}

	// 查询总记录数
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "count books failed: %v", err)
	}

	// 如果没有记录，直接返回空列表
	if total == 0 {
		return []*Book{}, 0, nil
	}

	// 计算分页的偏移量
	offset := int((req.PageNum - 1) * req.PageSize)
	if offset < 0 {
		offset = 0
	}

	// 查询书籍列表
	err = query.Order("id desc").
		Offset(offset).
		Limit(int(req.PageSize)).
		Find(&results).
		Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "search book failed: %v", err)
	}

	// 将结果转换为指针切片
	var resultBooks []*Book
	for i := range results {
		resultBooks = append(resultBooks, &results[i])
	}

	return resultBooks, total, nil
}

// IsBookExist 检查指定 ID 的书籍是否存在
// 1. 根据 BookID 查询书籍数量。
// 2. 如果数量大于 0，返回 true，否则返回 false。
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

// GetBookById 根据 ID 获取书籍信息
// 1. 根据 BookID 查询书籍。
// 2. 如果书籍存在，返回书籍信息，否则返回错误。
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

// IsBookInISBN 检查指定 ISBN 的书籍是否存在
// 1. 根据 ISBN 查询书籍数量。
// 2. 如果数量大于 0，返回 true，否则返回 false。
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
