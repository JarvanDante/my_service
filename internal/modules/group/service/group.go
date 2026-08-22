package service

import "context"

type ItemDTO struct {
	Id        int64
	Name      string
	Intro     string
	Avatar    string
	Link      string
	Platform  string
	Rank      int
	Status    int
	CreatedAt string
}

type SaveInput struct {
	Id       int64
	Name     string
	Intro    string
	Avatar   string
	Link     string
	Platform string
	Rank     int
	Status   int
}

type ListFilter struct {
	Status  int
	Keyword string
	Page    int
	Size    int
}

type IGroup interface {
	FrontList(ctx context.Context) ([]*ItemDTO, error)
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
}
