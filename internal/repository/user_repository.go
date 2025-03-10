package repository

import (
	"database/sql"
	"time"

	"github.com/daioru/todo-app/internal/logger"
	"github.com/daioru/todo-app/internal/models"
	"github.com/rs/zerolog"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db  *sqlx.DB
	sq  squirrel.StatementBuilderType
	log zerolog.Logger
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:  db,
		sq:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log: logger.GetLogger(),
	}
}

func (r *UserRepository) UserExists(user *models.User) (bool, error) {
	query, args, err := r.sq.Select("COUNT(*)").
		From("users").
		Where(squirrel.Eq{"username": user.Username}).
		ToSql()
	if err != nil {
		r.log.Error().
			Object("user", user).
			Err(err).
			Msg("Failed to build UserExists query")
		return false, err
	}

	var count int
	err = r.db.Get(&count, query, args...)
	if err != nil {
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("UserExists DB execution error")
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query, args, err := r.sq.Insert("users").
		Columns("username", "password_hash", "created_at").
		Values(user.Username, user.PasswordHash, time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		r.log.Error().
			Object("user", user).
			Err(err).
			Msg("Failed to build CreateUser query")
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&user.ID)
	if err != nil {
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("UserExists DB execution error")
		return err
	}

	return err
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User

	query, args, err := r.sq.Select("id", "username", "password_hash", "created_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		r.log.Error().
			Int("user_id", id).
			Err(err).
			Msg("Failed to build GetUserByID query")
		return nil, err
	}

	err = r.db.Get(&user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("GetUserByID DB execution error")
		return nil, err
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
		r.log.Error().
			Str("username", username).
			Err(err).
			Msg("Failed to build GetUserByUsername query")
		return nil, err
	}

	err = r.db.Get(&user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("GetUserByUsername DB execution error")
		return nil, err
	}

	return &user, nil
}
