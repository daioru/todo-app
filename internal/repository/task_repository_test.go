package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := repository.NewTaskRepository(db)
	task := &models.Task{
		UserID:      1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(task.UserID, task.Title, task.Description, task.Status, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err = repo.CreateTask(task)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasksByUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := repository.NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "status", "created_at"}).
		AddRow(1, 1, "Test Task", "Test Description", "pending", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM tasks").
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := repo.GetTasksByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Test Task", tasks[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTask(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := repository.NewTaskRepository(db)

	updates := make(map[string]interface{})
	updates["id"] = 1
	updates["user_id"] = 1
	updates["title"] = "Updated task"
	updates["description"] = "Updated description"
	updates["status"] = "done"

	mock.ExpectExec("UPDATE tasks").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), updates["id"], updates["user_id"]).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateTask(updates)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
