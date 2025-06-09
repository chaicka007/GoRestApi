package handlers

import (
	"RestApi/database"
	"RestApi/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTasks godoc
// @Summary Получить список задач
// @Description Возвращает все задачи или фильтрует по статусу
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "Фильтр по статусу: pending, in_progress, completed"
// @Success 200 {array} models.Task
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [get]
func GetTasks(c *gin.Context) {
	status := c.Query("status")

	var rows *sql.Rows
	var err error

	if status != "" {
		if !models.IsValidStatus(status) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Недопустимый фильтр статуса"})
			return
		}
		rows, err = database.DB.Query("SELECT id, title, description, status FROM tasks WHERE status = $1", status)
	} else {
		rows, err = database.DB.Query("SELECT id, title, description, status FROM tasks")
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при получении задач"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при обработке данных"})
			return
		}
		tasks = append(tasks, t)
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID godoc
// @Summary Получить задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var t models.Task
	err := database.DB.QueryRow("SELECT id, title, description FROM tasks WHERE id = $1", id).
		Scan(&t.ID, &t.Title, &t.Description, &t.Status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Задача не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при получении задачи"})
		return
	}
	c.JSON(http.StatusOK, t)
}

// CreateTask godoc
// @Summary Создать новую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Новая задача"
// @Success 201 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Неверный формат JSON"})
		return
	}
	if t.Title == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Заголовок задачи не может быть пустым"})
		return
	}
	if t.Status == "" {
		t.Status = models.StatusPending
	} else if !models.IsValidStatus(t.Status) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Недопустимый статус задачи"})
		return
	}
	err := database.DB.QueryRow("INSERT INTO tasks (title, description,status) VALUES ($1, $2, $3) RETURNING id",
		t.Title, t.Description, t.Status).Scan(&t.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при создании задачи"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// UpdateTask godoc
// @Summary Обновить задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body models.Task true "Обновлённая задача"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Неверный формат JSON"})
		return
	}
	if !models.IsValidStatus(t.Status) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Недопустимый статус задачи"})
		return
	}

	res, err := database.DB.Exec("UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
		t.Title, t.Description, t.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при обновлении задачи"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Задача не найдена"})
		return
	}
	t.ID, _ = strconv.Atoi(id)
	c.JSON(http.StatusOK, t)
}

// DeleteTask godoc
// @Summary Удалить задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	res, err := database.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при удалении задачи"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Задача не найдена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Задача удалена"})
}
