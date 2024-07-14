package controllers

import (
	"context"
	"fmt"
	"net/http"
	"simple-todo/config"
	"simple-todo/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var todoCollection *mongo.Collection = config.OpenCollection(config.Client, "todos")

func GetTodos(c *gin.Context) {
	// สร้าง context สำหรับการดึงข้อมูล
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"FindError": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	for cursor.Next(ctx) {
		var todo models.Todo

		err := cursor.Decode(&todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Decode Error": err.Error()})
			return
		}

		todos = append(todos, todo)
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"todos":   todos,
	})
}

func CreateTodo(c *gin.Context) {
	var todo models.Todo

	// ผูกข้อมูลจาก request body กับ struct User
	if err := c.ShouldBindJSON(&todo); err != nil {
		// ถ้ามีข้อผิดพลาดในการผูกข้อมูล
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.Id = primitive.NewObjectID()
	todo.Status = "todo"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := todoCollection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"todo":    todo,
	})
}

func GetTodoById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	var todo models.Todo
	err := todoCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"todo":    todo,
	})
}

func DeleteTodo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	deleteResult, err := todoCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if deleteResult.DeletedCount == 0 {
		msg := fmt.Sprintf("No todo with id : %v was found, no deletion occurred.", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"success": true,
		"id":      objId,
	})
}

func UpdateTodo(c *gin.Context) {
	var validate *validator.Validate = validator.New()

	var reqTodo models.UpdateTodo

	if err := c.BindJSON(&reqTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reqTodo.Task != nil && *reqTodo.Task == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task cannot be empty"})
		return
	}

	updateData := bson.M{}
	if reqTodo.Task != nil {
		updateData["task"] = *reqTodo.Task
	}

	if reqTodo.Status != nil {
		if err := validate.Var(reqTodo.Status, "oneof=todo doing done"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}
		updateData["status"] = *reqTodo.Status
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := todoCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateData})

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"todo":    reqTodo,
	})
}
