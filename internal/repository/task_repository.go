package repository

import (
	"database/sql"
	"time"

	"github.com/daioru/todo-app/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *TaskRepository) CreateTask(task *models.Task) error {
	query, args, err := r.sq.Insert("tasks").
		Columns("user_id", "title", "description", "status", "created_at").
		Values(task.UserID, task.Title, task.Description, task.Status, time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}
	return r.db.QueryRow(query, args...).Scan(&task.ID)
}

func (r *TaskRepository) GetTaskByID(id int) (*models.Task, error) {
	var task models.Task
	query, args, err := r.sq.Select("id", "user_id", "title", "description", "status", "created_at").
		From("tasks").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&task, query, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

func (r *TaskRepository) GetTasksByUserID(userID int) ([]models.Task, error) {
	var tasks []models.Task
	query, args, err := r.sq.Select("id", "user_id", "title", "description", "status", "created_at").
		From("tasks").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&tasks, query, args...)
	return tasks, err
}

func (r *TaskRepository) DeleteTask(taskID, userID int) error {
	query, args, err := r.sq.Delete("tasks").
		Where(squirrel.And{
			squirrel.Eq{"id": taskID},
			squirrel.Eq{"user_id": userID},
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *TaskRepository) UpdateTask(task *models.Task) error {
	query, args, err := r.sq.Update("tasks").
		Set("title", task.Title).
		Set("description", task.Description).
		Set("status", task.Status).
		Where(squirrel.And{
			squirrel.Eq{"id": task.ID},
			squirrel.Eq{"user_id": task.UserID},
		}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args)
	return err
}
