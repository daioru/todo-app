package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func NewMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *repository.UserRepository) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := repository.NewUserRepository(db)
	return mockDB, mock, repo
}

func TestCreateUser(t *testing.T) {
	mockDB, mock, repo := NewMock(t)
	defer mockDB.Close()

	user := &models.User{
		Username:     "test user",
		Password:     "test password",
		PasswordHash: "test password hash",
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Username, user.PasswordHash, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	mockDB, mock, repo := NewMock(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at"}).
		AddRow(1, "Test username", "Test password hash", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Test username", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername(t *testing.T) {
	mockDB, mock, repo := NewMock(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at"}).
		AddRow(1, "Test username", "Test password hash", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs("username").
		WillReturnRows(rows)

	user, err := repo.GetUserByUsername("username")
	assert.NoError(t, err)
	assert.Equal(t, "Test username", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}
