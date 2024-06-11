package types

import (
	"time"
)

type CreateNewNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewNote(Title string, Body string) *Note {
	return &Note{
		Title:     Title,
		Body:      Body,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
