package service

import (
	"context"
	"films-service/internal/model"
)

type FilmRepository interface {
	GetAll(c context.Context) ([]model.Film, error)
}

type FilmService struct {
	repo FilmRepository
}

func NewFilmService(repo FilmRepository) *FilmService {
	if repo == nil {
		panic("repo is nil")
	}
	return &FilmService{repo: repo}
}

func (s *FilmService) GetFilms(ctx context.Context) ([]model.Film, error) {
	return s.repo.GetAll(ctx)
}
