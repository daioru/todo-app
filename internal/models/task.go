package models

import (
	"time"

	"github.com/rs/zerolog"
)

type Task struct {
	ID          int      `db:"id" json:"id"`
	UserID      int      `db:"user_id" json:"user_id"`
	Title       string   `db:"title" json:"title" binding:"required"`
	Description string   `db:"description" json:"description" binding:"required"`
	Status      string   `db:"status" json:"status" binding:"required"`
	CreatedAt   JSONTime `db:"created_at" json:"created_at"`
}

func (t Task) MarshalZerologObject(e *zerolog.Event) {
	e.Int("id", t.ID).
		Int("user_id", t.UserID).
		Str("title", t.Title).
		Str("description", t.Description).
		Str("status", t.Status).
		Time("created_at", time.Time(t.CreatedAt))
}
