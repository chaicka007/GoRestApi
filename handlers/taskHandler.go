package handlers

import (
	"RestApi/models"
	"RestApi/storage"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Storage *storage.TaskStorage
}

func NewTaskHandler(storage *storage.TaskStorage) *TaskHandler {
	return &TaskHandler{Storage: storage}
}

// GetTasks godoc
// @Summary Получить список задач
// @Description Возвращает все задачи или фильтрует по статусу
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "Фильтр по статусу: pending, in_progress, completed"
// @Success 200 {array} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	status := c.Query("status")

	tasks, err := h.Storage.GetAll(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при получении задач"})
		return
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
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Некорректный ID"})
		return
	}

	task, err := h.Storage.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при получении задачи"})
		return
	}
	if task == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Задача не найдена"})
		return
	}
	c.JSON(http.StatusOK, task)
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
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Неверные данные"})
		return
	}

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Название задачи не может быть пустым"})
		return
	}

	if err := h.Storage.Create(&task); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при создании задачи"})
		return
	}
	c.JSON(http.StatusCreated, task)
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
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Некорректный ID"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Неверные данные"})
		return
	}

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Название задачи не может быть пустым"})
		return
	}

	if err := h.Storage.Update(id, &task); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Задача не найдена"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при обновлении задачи"})
		}
		return
	}

	task.ID = id
	c.JSON(http.StatusOK, task)
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
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Некорректный ID"})
		return
	}

	if err := h.Storage.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Задача не найдена"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Ошибка при удалении задачи"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
