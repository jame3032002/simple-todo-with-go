package main

import (
	"net/http"
	"simple-todo/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Server is running",
		})
	})

	r.GET("/todos", controllers.GetTodos)
	r.POST("/todos", controllers.CreateTodo)
	r.GET("/todos/:id", controllers.GetTodoById)
	r.DELETE("/todos/:id", controllers.DeleteTodo)
	r.PATCH("/todos/:id", controllers.UpdateTodo)
	r.Run(":2000")
}
