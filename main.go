package main

import (
	"RestApi/database"
	"RestApi/handlers"
	"RestApi/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "RestApi/docs"
)

// @title Task API
// @version 1.0
// @description Документация API
// @host localhost:8080
// @BasePath /

func main() {
	database.InitDB()
	r := gin.Default()
	h := handlers.NewTaskHandler(storage.NewTaskStorage())

	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/tasks", h.GetTasks)
	r.GET("/tasks/:id", h.GetTaskByID)
	r.POST("/tasks", h.CreateTask)
	r.PUT("/tasks/:id", h.UpdateTask)
	r.DELETE("/tasks/:id", h.DeleteTask)

	r.Run(":8080")
}
