package usecase

import (
	"context"

	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/entity"
)

type (
	Movies interface {
		Index(ctx context.Context, data []entity.Movie) error
		Autocomplete(ctx context.Context, query string) ([]entity.Movie, error)
		Search(ctx context.Context, query string) ([]entity.Movie, error)
	}

	MoviesRepo interface {
		Index(ctx context.Context, data []entity.Movie) error
		Autocomplete(ctx context.Context, query string) ([]entity.Movie, error)
		Search(ctx context.Context, query string) ([]entity.Movie, error)
		CreateIndex(ctx context.Context) error
	}
)
