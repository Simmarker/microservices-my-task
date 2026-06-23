package main

import (
	"context"
	"films-service/internal/handler"
	"films-service/internal/repository"
	"films-service/internal/service"
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

//TODO: http вместо gin

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conn, err := pgx.Connect(ctx,
		"postgres://program:test@localhost:5432/films")
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer conn.Close(ctx)

	repo := repository.NewFilmRepository(conn)
	svc := service.NewFilmService(repo)
	h := handler.NewFilmHandler(svc)

	router := gin.Default()

	router.GET("/api/v1/films", h.GetFilms)

	if err := router.Run(":8070"); err != nil {
		log.Fatalf("Ошибка запуска: %v", err)
	}
}
