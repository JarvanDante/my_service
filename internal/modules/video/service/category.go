package service

import "context"

type CategoryDTO struct {
	Id        int64
	Name      string
	Kind      int
	Rank      int
	Status    int
	CreatedAt string
}

type CategoryInput struct {
	Id     int64
	Name   string
	Kind   int
	Rank   int
	Status int
}

type CategoryFilter struct {
	Kind   int
	Status int
	Page   int
	Size   int
}

type ICategory interface {
	List(ctx context.Context, f CategoryFilter) ([]*CategoryDTO, int, error)
	Create(ctx context.Context, in CategoryInput) (int64, error)
	Update(ctx context.Context, in CategoryInput) error
	Delete(ctx context.Context, id int64) error
}
