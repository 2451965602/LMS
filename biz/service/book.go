package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/book"
)

type BookService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewBookService(ctx context.Context, c *app.RequestContext) *BookService {
	return &BookService{
		ctx: ctx,
		c:   c,
	}
}

func (s *BookService) AddBook(ctx context.Context, req book.AddBookRequest) (int64, error) {
	bookId, err := db.AddBook(ctx, req)
	if err != nil {
		return -1, err
	}

	return bookId, nil
}

func (s *BookService) UpdateBook(ctx context.Context, req book.UpdateBookRequest) (*db.Book, error) {
	bk, err := db.UpdateBook(ctx, req)
	if err != nil {
		return nil, err
	}
	return bk, nil
}

func (s *BookService) DeleteBook(ctx context.Context, req book.DeleteBookRequest) error {
	err := db.DeleteBook(ctx, req.BookID)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookService) SearchBook(ctx context.Context, req book.GetBookRequest) ([]*db.Book, int64, error) {
	books, total, err := db.SearchBook(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (s *BookService) GetBookById(ctx context.Context, bookId int64) (*db.Book, error) {
	bk, err := db.GetBookById(ctx, bookId)
	if err != nil {
		return nil, err
	}
	return bk, nil
}
