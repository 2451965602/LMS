package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/book"
)

// BookService 用于管理图书相关的业务逻辑，封装了添加、更新、删除和查询图书的操作。
type BookService struct {
	ctx context.Context     // 上下文，用于传递请求相关的元数据
	c   *app.RequestContext // Hertz框架的请求上下文，用于处理HTTP请求
}

// NewBookService 创建一个新的BookService实例，初始化上下文和请求上下文。
func NewBookService(ctx context.Context, c *app.RequestContext) *BookService {
	return &BookService{
		ctx: ctx,
		c:   c,
	}
}

// AddBook 添加新的图书
// 参数：
//   - ctx: 上下文
//   - req: 添加图书请求，包含图书信息
//
// 返回值：
//   - int64: 添加成功的图书ID
//   - error: 错误信息，如果添加失败会返回错误
func (s *BookService) AddBook(ctx context.Context, req book.AddBookRequest) (int64, error) {
	bookId, err := db.AddBook(ctx, req) // 调用数据库操作函数添加图书
	if err != nil {
		return -1, err
	}

	return bookId, nil
}

// UpdateBook 更新图书信息
// 参数：
//   - ctx: 上下文
//   - req: 更新图书请求，包含图书信息
//
// 返回值：
//   - *db.Book: 更新成功的图书信息
//   - error: 错误信息，如果更新失败会返回错误
func (s *BookService) UpdateBook(ctx context.Context, req book.UpdateBookRequest) (*db.Book, error) {
	bk, err := db.UpdateBook(ctx, req) // 调用数据库操作函数更新图书
	if err != nil {
		return nil, err
	}
	return bk, nil
}

// DeleteBook 删除图书
// 参数：
//   - ctx: 上下文
//   - req: 删除图书请求，包含图书ID
//
// 返回值：
//   - error: 错误信息，如果删除失败会返回错误
func (s *BookService) DeleteBook(ctx context.Context, req book.DeleteBookRequest) error {
	err := db.DeleteBook(ctx, req.BookID) // 调用数据库操作函数删除图书
	if err != nil {
		return err
	}
	return nil
}

// SearchBook 搜索图书
// 参数：
//   - ctx: 上下文
//   - req: 搜索图书请求，包含搜索条件和分页信息
//
// 返回值：
//   - []*db.Book: 搜索结果的图书列表
//   - int64: 总记录数
//   - error: 错误信息，如果搜索失败会返回错误
func (s *BookService) SearchBook(ctx context.Context, req book.GetBookRequest) ([]*db.Book, int64, error) {
	books, total, err := db.SearchBook(ctx, req) // 调用数据库操作函数搜索图书
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// GetBookById 根据图书ID获取图书信息
// 参数：
//   - ctx: 上下文
//   - bookId: 图书ID
//
// 返回值：
//   - *db.Book: 获取的图书信息
//   - error: 错误信息，如果获取失败会返回错误
func (s *BookService) GetBookById(ctx context.Context, bookId int64) (*db.Book, error) {
	bk, err := db.GetBookById(ctx, bookId) // 调用数据库操作函数根据ID获取图书信息
	if err != nil {
		return nil, err
	}
	return bk, nil
}
