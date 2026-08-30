package service

import "context"

type ItemDTO struct {
	Id        int64
	Position  string
	Title     string
	CoverUrl  string
	Link      string
	Rank      int
	Status    int
	CreatedAt string
}

type SaveInput struct {
	Id       int64
	Position string
	Title    string
	CoverUrl string
	Link     string
	Rank     int
	Status   int
}

type ListFilter struct {
	Position string
	Status   int
	Keyword  string
	Page     int
	Size     int
}

type IBanner interface {
	FrontList(ctx context.Context, position string) ([]*ItemDTO, error)
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
}
