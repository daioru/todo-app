package repository

import (
	"database/sql"
	"time"

	"github.com/daioru/todo-app/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query, args, err := r.sq.Insert("users").
		Columns("username", "password_hash", "created_at").
		Values(user.Username, user.PasswordHash, time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}
	return r.db.QueryRow(query, args...).Scan(&user.ID)
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query, args, err := r.sq.Select("id", "username", "password_hash", "created_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&user, query, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query, args, err := r.sq.Select("id", "username", "password_hash", "created_at").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&user, query, args...)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, err
}
