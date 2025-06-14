package usecase

import (
	"context"

	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/entity"
)

type moviesUsecase struct {
	repo MoviesRepo
}

func NewMovies(repo MoviesRepo) Movies {
	return &moviesUsecase{repo: repo}
}

func (u *moviesUsecase) Index(ctx context.Context, data []entity.Movie) error {
	return u.repo.Index(ctx, data)
}

func (u *moviesUsecase) Autocomplete(ctx context.Context, query string) ([]entity.Movie, error) {
	if len(query) < 2 {
		return nil, nil
	}
	return u.repo.Autocomplete(ctx, query)
}

func (u *moviesUsecase) Search(ctx context.Context, query string) ([]entity.Movie, error) {
	if len(query) < 2 {
		return nil, nil
	}

	return u.repo.Search(ctx, query)
}
