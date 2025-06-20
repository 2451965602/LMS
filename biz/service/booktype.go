package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/booktype"
	"github.com/2451965602/LMS/pkg/errno"
)

// BookTypeService 用于管理图书类型相关的业务逻辑，封装了添加、更新、删除和查询图书类型的操作。
type BookTypeService struct {
	ctx context.Context     // 上下文，用于传递请求相关的元数据
	c   *app.RequestContext // Hertz框架的请求上下文，用于处理HTTP请求
}

// NewBookTypeService 创建一个新的BookTypeService实例，初始化上下文和请求上下文。
func NewBookTypeService(ctx context.Context, c *app.RequestContext) *BookTypeService {
	return &BookTypeService{
		ctx: ctx,
		c:   c,
	}
}

// AddBookType 添加新的图书类型
// 参数：
//   - ctx: 上下文
//   - req: 添加图书类型请求，包含图书类型信息
//
// 返回值：
//   - *db.BookType: 添加成功的图书类型信息
//   - error: 错误信息，如果添加失败会返回错误
func (s *BookTypeService) AddBookType(ctx context.Context, req booktype.AddBookTypeRequest) (*db.BookType, error) {

	// 检查ISBN格式是否正确
	if !IsValidISBN(req.ISBN) {
		return nil, errno.Errorf(errno.ServiceInvalidISBN, "invalid ISBN format") // 如果ISBN格式不正确，返回错误
	}

	if !CheckAuthor(req.Author) {
		return nil, errno.Errorf(errno.ServiceInvalidAuthor, "invalid author format") // 如果作者格式不正确，返回错误
	}

	exit, err := db.IsBookTypeExist(ctx, req.ISBN) // 检查图书类型是否已存在
	if err != nil {
		return nil, err
	}
	if exit {
		return nil, errno.Errorf(errno.ServiceBookTypeExist, "book type already exist") // 如果图书类型已存在，返回错误
	}

	bt, err := db.AddBookType(ctx, req) // 调用数据库操作函数添加图书类型
	if err != nil {
		return nil, err
	}
	return bt, nil
}

// UpdateBookType 更新图书类型信息
// 参数：
//   - ctx: 上下文
//   - req: 更新图书类型请求，包含图书类型信息
//
// 返回值：
//   - *db.BookType: 更新成功的图书类型信息
//   - error: 错误信息，如果更新失败会返回错误
func (s *BookTypeService) UpdateBookType(ctx context.Context, req booktype.UpdateBookTypeRequest) (*db.BookType, error) {

	// 检查ISBN格式是否正确
	if !IsValidISBN(req.ISBN) {
		return nil, errno.Errorf(errno.ServiceInvalidISBN, "invalid ISBN format") // 如果ISBN格式不正确，返回错误
	}
	if req.Author != nil {
		if !CheckAuthor(*req.Author) {
			return nil, errno.Errorf(errno.ServiceInvalidAuthor, "invalid author format") // 如果作者格式不正确，返回错误
		}
	}
	bt, err := db.UpdateBookType(ctx, req) // 调用数据库操作函数更新图书类型
	if err != nil {
		return nil, err
	}
	return bt, nil
}

// DeleteBookType 删除图书类型
// 参数：
//   - ctx: 上下文
//   - req: 删除图书类型请求，包含ISBN
//
// 返回值：
//   - error: 错误信息，如果删除失败会返回错误
func (s *BookTypeService) DeleteBookType(ctx context.Context, req booktype.DeleteBookTypeRequest) error {
	// 检查ISBN格式是否正确
	if !IsValidISBN(req.ISBN) {
		return errno.Errorf(errno.ServiceInvalidISBN, "invalid ISBN format") // 如果ISBN格式不正确，返回错误
	}

	exist, err := db.IsBookTypeExist(ctx, req.ISBN) // 检查图书类型是否存在
	if err != nil {
		return err
	}
	if !exist {
		return errno.Errorf(errno.ServiceBookTypeNotExist, "book type not exist") // 如果图书类型不存在，返回错误
	}

	exit, err := db.IsBookInISBN(ctx, req.ISBN) // 检查是否有书籍使用该图书类型
	if err != nil {
		return err
	}
	if exit {
		return errno.Errorf(errno.ServiceBookTypeInUse, "book type cannot be deleted, existing books of this type still exist") // 如果有书籍使用该类型，返回错误
	}

	err = db.DeleteBookType(ctx, req.ISBN) // 调用数据库操作函数删除图书类型
	if err != nil {
		return err
	}
	return nil
}

// SearchBookType 搜索图书类型
// 参数：
//   - ctx: 上下文
//   - req: 搜索图书类型请求，包含搜索条件和分页信息
//
// 返回值：
//   - []*db.BookType: 搜索结果的图书类型列表
//   - int64: 总记录数
//   - error: 错误信息，如果搜索失败会返回错误
func (s *BookTypeService) SearchBookType(ctx context.Context, req booktype.GetBookTypeRequest) ([]*db.BookType, int64, error) {

	if req.ISBN != nil {
		// 检查ISBN格式是否正确
		if !IsValidISBN(*req.ISBN) {
			return nil, 0, errno.Errorf(errno.ServiceInvalidISBN, "invalid ISBN format") // 如果ISBN格式不正确，返回错误
		}
	}

	bookTypes, total, err := db.SearchBookType(ctx, req.Title, req.Author, req.ISBN, req.PageNum, req.PageSize) // 调用数据库操作函数搜索图书类型
	if err != nil {
		return nil, 0, err
	}

	return bookTypes, total, nil
}

// GetBookTypeByISBN 根据ISBN获取图书类型
// 参数：
//   - ctx: 上下文
//   - isbn: 图书类型的ISBN
//
// 返回值：
//   - *db.BookType: 获取的图书类型信息
//   - error: 错误信息，如果获取失败会返回错误
func (s *BookTypeService) GetBookTypeByISBN(ctx context.Context, isbn string) (*db.BookType, error) {
	// 检查ISBN格式是否正确
	if !IsValidISBN(isbn) {
		return nil, errno.Errorf(errno.ServiceInvalidISBN, "invalid ISBN format") // 如果ISBN格式不正确，返回错误
	}
	bt, err := db.GetBookTypeByISBN(ctx, isbn) // 调用数据库操作函数根据ISBN获取图书类型
	if err != nil {
		return nil, err
	}
	return bt, nil
}
