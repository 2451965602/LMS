package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/2451965602/LMS/biz/dal/db"
	"github.com/2451965602/LMS/biz/model/booktype"
	"github.com/2451965602/LMS/pkg/errno"
)

type BookTypeService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewBookTypeService(ctx context.Context, c *app.RequestContext) *BookTypeService {
	return &BookTypeService{
		ctx: ctx,
		c:   c,
	}
}

func (s *BookTypeService) AddBookType(ctx context.Context, req booktype.AddBookTypeRequest) (*db.BookType, error) {
	exit, err := db.IsBookTypeExist(ctx, req.ISBN)
	if err != nil {
		return nil, err
	}
	if exit {
		return nil, errno.Errorf(errno.ServiceBookTypeExist, "book type already exist")
	}

	bt, err := db.AddBookType(ctx, req)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func (s *BookTypeService) UpdateBookType(ctx context.Context, req booktype.UpdateBookTypeRequest) (*db.BookType, error) {
	bt, err := db.UpdateBookType(ctx, req)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func (s *BookTypeService) DeleteBookType(ctx context.Context, req booktype.DeleteBookTypeRequest) error {
	exist, err := db.IsBookTypeExist(ctx, req.ISBN)
	if err != nil {
		return err
	}
	if !exist {
		return errno.Errorf(errno.ServiceBookTypeNotExist, "book type not exist")
	}

	exit, err := db.IsBookInISBN(ctx, req.ISBN)
	if err != nil {
		return err
	}
	if exit {
		return errno.Errorf(errno.ServiceBookTypeInUse, "book type cannot be deleted, existing books of this type still exist")
	}

	err = db.DeleteBookType(ctx, req.ISBN)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookTypeService) SearchBookType(ctx context.Context, req booktype.GetBookTypeRequest) ([]*db.BookType, int64, error) {
	bookTypes, total, err := db.SearchBookType(ctx, req.Title, req.ISBN, req.Author, req.PageNum, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	return bookTypes, total, nil
}

func (s *BookTypeService) GetBookTypeByISBN(ctx context.Context, isbn string) (*db.BookType, error) {
	bt, err := db.GetBookTypeByISBN(ctx, isbn)
	if err != nil {
		return nil, err
	}
	return bt, nil
}
