package service

import (
	"context"
	"github.com/2451965602/LMS/config"
	"github.com/2451965602/LMS/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/borrow"
	contextLogin "github.com/2451965602/LMS/pkg/base/context"
)

// BorrowService 用于管理图书借阅相关的业务逻辑，封装了借书、还书、续借和获取借阅记录等操作。
type BorrowService struct {
	ctx context.Context     // 上下文，用于传递请求相关的元数据
	c   *app.RequestContext // Hertz框架的请求上下文，用于处理HTTP请求
}

// NewBorrowService 创建一个新的BorrowService实例，初始化上下文和请求上下文。
func NewBorrowService(ctx context.Context, c *app.RequestContext) *BorrowService {
	return &BorrowService{
		ctx: ctx,
		c:   c,
	}
}

// BookBorrow 借书操作
// 参数：
//   - ctx: 上下文
//   - req: 借书请求，包含书籍ID
//
// 返回值：
//   - int64: 借阅记录ID
//   - error: 错误信息，如果借书失败会返回错误
func (s *BorrowService) BookBorrow(ctx context.Context, req borrow.BorrowRequest) (int64, error) {
	userId, err := contextLogin.GetLoginData(ctx) // 从上下文中获取当前登录用户ID
	if err != nil {
		return -1, err
	}
	_, count, err := db.GetCurrentBorrowRecord(ctx, userId, 1, 1, 1)
	if err != nil {
		return -1, err
	}
	if count >= config.MaxBorrowNum.Num {
		return -1, errno.Errorf(errno.ServiceBorrowNumOver, "can't borrow more")
	}

	borrowId, err := db.BookBorrow(ctx, userId, req.BookID) // 调用数据库操作函数记录借书信息
	if err != nil {
		return -1, err
	}
	return borrowId, nil
}

// BookReturn 还书操作
// 参数：
//   - ctx: 上下文
//   - req: 还书请求，包含书籍ID、借阅记录ID、还书状态和逾期费用
//
// 返回值：
//   - *db.BorrowRecord: 还书后的借阅记录信息
//   - error: 错误信息，如果还书失败会返回错误
func (s *BorrowService) BookReturn(ctx context.Context, req borrow.ReturnRequest) (*db.BorrowRecord, error) {
	userId, err := contextLogin.GetLoginData(ctx) // 从上下文中获取当前登录用户ID
	if err != nil {
		return nil, err
	}
	borrowRecord, err := db.BookReturn(ctx, userId, req.BookID, req.BorrowID, req.Status, req.LateFee) // 调用数据库操作函数记录还书信息
	if err != nil {
		return nil, err
	}
	return borrowRecord, nil
}

// BookRenew 续借操作
// 参数：
//   - ctx: 上下文
//   - req: 续借请求，包含借阅记录ID和额外借阅时间
//
// 返回值：
//   - *db.BorrowRecord: 续借后的借阅记录信息
//   - error: 错误信息，如果续借失败会返回错误
func (s *BorrowService) BookRenew(ctx context.Context, req borrow.RenewRequest) (*db.BorrowRecord, error) {
	userId, err := contextLogin.GetLoginData(ctx) // 从上下文中获取当前登录用户ID
	if err != nil {
		return nil, err
	}
	borrowRecord, err := db.BookRenew(ctx, userId, req.BorrowID, int(req.AddTime)) // 调用数据库操作函数记录续借信息
	if err != nil {
		return nil, err
	}
	return borrowRecord, nil
}

// GetCurrentBorrowRecord 获取当前用户的借阅记录
// 参数：
//   - ctx: 上下文
//   - req: 获取借阅记录请求，包含分页信息和借阅状态
//
// 返回值：
//   - []*db.BorrowRecord: 借阅记录列表
//   - int64: 总记录数
//   - error: 错误信息，如果获取失败会返回错误
func (s *BorrowService) GetCurrentBorrowRecord(ctx context.Context, req borrow.GetBorrowRecordRequest) ([]*db.BorrowRecord, int64, error) {
	userId, err := contextLogin.GetLoginData(ctx) // 从上下文中获取当前登录用户ID
	if err != nil {
		return nil, 0, err
	}

	records, total, err := db.GetCurrentBorrowRecord(ctx, userId, req.PageNum, req.PageSize, req.Status) // 调用数据库操作函数获取借阅记录
	if err != nil {
		return nil, 0, err
	}

	var resultRecords []*db.BorrowRecord
	for i := range records {
		resultRecords = append(resultRecords, &records[i]) // 将借阅记录转换为指针列表
	}
	return resultRecords, total, nil
}
