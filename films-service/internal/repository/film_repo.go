package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"films-service/internal/model"
)

type FilmRepository struct {
	conn *pgx.Conn
}

func NewFilmRepository(conn *pgx.Conn) *FilmRepository {
	return &FilmRepository{conn: conn}
}

//TODO: sqlc вместо Query

func (r *FilmRepository) GetAll(ctx context.Context) ([]model.Film, error) {
	rows, err := r.conn.Query(ctx, "SELECT film_uid, name, rating, director, producer, genre FROM film")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []model.Film
	for rows.Next() {
		var f model.Film
		if err := rows.Scan(&f.FilmUID, &f.Name, &f.Rating, &f.Director, &f.Producer, &f.Genre); err != nil {
			return nil, err
		}
		films = append(films, f)
	}
	return films, rows.Err()
}
