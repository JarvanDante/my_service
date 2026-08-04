package service

import "context"

type VideoDTO struct {
	Id            int64
	Title         string
	Description   string
	CoverUrl      string
	CoverKey      string
	CoverMediaId  int64
	SourceUrl     string
	SourceKey     string
	SourceMediaId int64
	Duration      int
	Sort          int
	Status        int
	CreatedBy     int64
	CreatedAt     string
	UpdatedAt     string
}

type ListInput struct {
	Keyword string
	Status  int
	Page    int
	Size    int
}

type ListDTO struct {
	List  []*VideoDTO
	Total int
	Page  int
	Size  int
}

type SaveInput struct {
	Id            int64
	Title         string
	Description   string
	CoverUrl      string
	CoverKey      string
	CoverMediaId  int64
	SourceUrl     string
	SourceKey     string
	SourceMediaId int64
	Duration      int
	Sort          int
	Status        int
	OperatorId    int64
}

type IVideo interface {
	List(ctx context.Context, in ListInput) (*ListDTO, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
	SetStatus(ctx context.Context, id int64, status int) error
}
