package service

import "context"

type ItemDTO struct {
	Id             int64
	Code           string
	Name           string
	Remark         string
	Status         int
	StatusText     string
	UserCount      int
	H5Link         string
	ClipboardValue string
	CreatedAt      string
}

type SaveInput struct {
	Id     int64
	Code   string
	Name   string
	Remark string
	Status int
}

type ListFilter struct {
	Status  int
	Keyword string
	Page    int
	Size    int
}

type IChannel interface {
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
}
