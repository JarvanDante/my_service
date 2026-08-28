package service

import "context"

type ModuleDTO struct {
	Id        int64
	Name      string
	Position  string
	Style     int
	Icon      int
	TagIds    []int64
	TagNames  []string
	Size      int
	Rank      int
	Status    int
	CreatedAt string
	UpdatedAt string
}

type ModuleInput struct {
	Id       int64
	Name     string
	Position string
	Style    int
	Icon     int
	TagIds   []int64
	Size     int
	Rank     int
	Status   int
}

type ModuleFilter struct {
	Name     string
	Position string
	Status   int // -1=全部
	Page     int
	Size     int
}

type ModuleFrontDTO struct {
	Id    int64
	Name  string
	Style int
	Icon  int
	Size  int
	Tags  []string
	Items []*ComicsDTO
}

type IModule interface {
	List(ctx context.Context, f ModuleFilter) ([]*ModuleDTO, int, error)
	Create(ctx context.Context, in ModuleInput) (int64, error)
	Update(ctx context.Context, in ModuleInput) error
	Delete(ctx context.Context, id int64) error
	FrontRepo(ctx context.Context, position string) ([]*ModuleFrontDTO, error)
}
