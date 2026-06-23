package model

import "github.com/google/uuid"

type Film struct {
	FilmUID  uuid.UUID `json:"film_uid"`
	Name     string    `json:"name"`
	Rating   int       `json:"rating"`
	Director string    `json:"director"`
	Producer string    `json:"producer"`
	Genre    string    `json:"genre"`
}
