package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"films-service/internal/model"
)

//TODO: http вместо gin

type FilmService interface {
	GetFilms(c context.Context) ([]model.Film, error)
}

type FilmHandler struct {
	service FilmService
}

func NewFilmHandler(service FilmService) *FilmHandler {
	return &FilmHandler{service: service}
}

func (h *FilmHandler) GetFilms(c *gin.Context) {
	films, err := h.service.GetFilms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения фильмов"})
		return
	}
	c.JSON(http.StatusOK, films)
}
