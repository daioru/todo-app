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

type TaskRepository struct {
	db  *sqlx.DB
	sq  squirrel.StatementBuilderType
	log zerolog.Logger
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		db:  db,
		sq:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log: logger.GetLogger(),
	}
}

func (r *TaskRepository) CreateTask(task *models.Task) error {
	query, args, err := r.sq.Insert("tasks").
		Columns("user_id", "title", "description", "status", "created_at").
		Values(task.UserID, task.Title, task.Description, task.Status, time.Now()).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		r.log.Error().
			Object("task", task).
			Err(err).
			Msg("Failed to build CreateTask query")
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("CreateTask DB execution error")
		return err
	}

	return nil
}

func (r *TaskRepository) GetTaskByID(id int) (*models.Task, error) {
	var task models.Task

	query, args, err := r.sq.Select("id", "user_id", "title", "description", "status", "created_at").
		From("tasks").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		r.log.Error().
			Int("task_id", id).
			Err(err).
			Msg("Failed to build GetTaskByID query")
		return &task, err
	}

	err = r.db.Get(&task, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("GetTaskByID DB execution error")
		return &task, err
	}

	return &task, nil
}

func (r *TaskRepository) GetTasksByUserID(userID int) ([]models.Task, error) {
	var tasks []models.Task

	query, args, err := r.sq.Select("id", "user_id", "title", "description", "status", "created_at").
		From("tasks").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		r.log.Error().
			Int("user_id", userID).
			Err(err).
			Msg("Failed to build GetTasksByUserID query")
		return tasks, err
	}

	err = r.db.Select(&tasks, query, args...)
	if err != nil {
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("GetTasksByUserID DB execution error")
		return tasks, err
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(taskID, userID int) error {
	query, args, err := r.sq.Delete("tasks").
		Where(squirrel.And{
			squirrel.Eq{"id": taskID},
			squirrel.Eq{"user_id": userID},
		}).
		ToSql()
	if err != nil {
		r.log.Error().
			Int("task_id", taskID).
			Int("user_id", userID).
			Err(err).
			Msg("Failed to build DeleteTask query")
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("DeleteTask DB execution error")
		return err
	}

	return nil
}

func (r *TaskRepository) UpdateTask(task *models.Task) error {
	query, args, err := r.sq.Update("tasks").
		Set("title", task.Title).
		Set("description", task.Description).
		Set("status", task.Status).
		Where(squirrel.Eq{"id": task.ID, "user_id": task.UserID}).
		Suffix("RETURNING created_at").
		ToSql()
	if err != nil {
		r.log.Error().
			Int("task_id", task.ID).
			Int("user_id", task.UserID).
			Err(err).
			Msg("Failed to build UpdateTask query")
		return err
	}

	err = r.db.Get(&task.CreatedAt, query, args...)
	if err != nil {
		r.log.Error().
			Str("query", query).
			Interface("args", args).
			Err(err).
			Msg("UpdateTask DB execution error")
		return err
	}

	return nil
}
