package service

import "context"

type ItemDTO struct {
	Id           int64
	Name         string
	IconUrl      string
	OpenMode     string
	OpenModeName string
	Link         string
	AppLink      string
	LinkLabel    string
	Position     string
	PositionName string
	Sort         int
	Status       int
	StatusText   string
	CreatedAt    string
	UpdatedAt    string
}

type SaveInput struct {
	Id       int64
	Name     string
	IconUrl  string
	OpenMode string
	Link     string
	AppLink  string
	Position string
	Sort     int
	Status   int
}

type ListFilter struct {
	Name     string
	Position string
	Status   int // -1=全部
	Page     int
	Size     int
}

type IKingkong interface {
	FrontList(ctx context.Context, position string) ([]*ItemDTO, error)
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
}
