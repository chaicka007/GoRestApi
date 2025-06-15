package storage

import (
	"RestApi/database"
	"RestApi/models"
	"database/sql"
)

type TaskStorage struct {
	DB *sql.DB
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		DB: database.DB,
	}
}

func (s *TaskStorage) GetAll(status string) ([]models.Task, error) {
	var rows *sql.Rows
	var err error

	if status != "" {
		rows, err = s.DB.Query("SELECT id, title, description, status FROM tasks WHERE status = $1", status)
	} else {
		rows, err = s.DB.Query("SELECT id, title, description, status FROM tasks")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *TaskStorage) GetByID(id int) (*models.Task, error) {
	row := s.DB.QueryRow("SELECT id, title, description, status FROM tasks WHERE id = $1", id)
	var t models.Task
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (s *TaskStorage) Create(task *models.Task) error {
	return s.DB.QueryRow(
		"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id",
		task.Title, task.Description, task.Status,
	).Scan(&task.ID)
}

func (s *TaskStorage) Update(id int, task *models.Task) error {
	result, err := s.DB.Exec(
		"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
		task.Title, task.Description, task.Status, id,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *TaskStorage) Delete(id int) error {
	result, err := s.DB.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
