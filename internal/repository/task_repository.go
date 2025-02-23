package repository

import (
	"database/sql"
	"time"

	"github.com/daioru/todo-app/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

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
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		log.WithFields(logrus.Fields{
			"user_id":     task.UserID,
			"task_id":     task.UserID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
		}).Errorf("Failed to build SQL query: %v", err)
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
		}).Errorf("DB execution error: %v", err)
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
		log.WithFields(logrus.Fields{
			"task_id": id,
		}).Errorf("Failed to build SQL query: %v", err)
		return &task, err
	}

	err = r.db.Get(&task, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
		}).Errorf("DB execution error: %v", err)
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
		log.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("Failed to build SQL query: %v", err)
		return tasks, err
	}

	err = r.db.Select(&tasks, query, args...)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
		}).Errorf("DB execution error: %v", err)
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
		log.WithFields(logrus.Fields{
			"task_id": taskID,
			"user_id": userID,
		}).Errorf("Failed to build SQL query: %v", err)
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
		}).Errorf("DB execution error: %v", err)
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
		log.WithFields(logrus.Fields{
			"task_id": task.ID,
			"user_id": task.UserID,
		}).Errorf("Failed to build SQL query: %v", err)
		return err
	}

	err = r.db.Get(&task.CreatedAt, query, args...)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
		}).Errorf("DB execution error: %v", err)
		return err
	}

	return nil
}
