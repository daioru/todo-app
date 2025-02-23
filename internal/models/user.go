package models

import (
	"time"

	"github.com/rs/zerolog"
)

type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Password     string    `db:"-" json:"password"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func (u User) MarshalZerologObject(e *zerolog.Event) {
	e.Int("id", u.ID).
		Str("username", u.Username).
		Str("password", u.Password).
		Str("password_hash", u.PasswordHash).
		Time("created_at", u.CreatedAt)
}
