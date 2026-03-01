package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Learn Go", Completed: true},
	{ID: "2", Item: "Build a web app", Completed: false},
	{ID: "3", Item: "Deploy to production", Completed: false},
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func addTodo(c *gin.Context) {
	var newTodo todo

	if err := c.BindJSON(&newTodo); err != nil {
		return
	}
    
	// check if the id already exists
	for _, t := range todos {
		if t.ID == newTodo.ID {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "todo with this ID already exists"})
			return
		}
	}
	
	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

// getTodo locates the todo whose ID value matches the id parameter sent by the client, then returns that todo as a response.
func getTodo(c *gin.Context) {
	id := c.Param("id")

	foundTodo, err := getTodoByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, foundTodo) // 200 OK
}

func getTodoByID(id string) (*todo, error) {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("todo not found")
}

func toggleTodoStatus(context *gin.Context) {
	id := context.Param("id")

	foundTodo, err := getTodoByID(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}

	foundTodo.Completed = !foundTodo.Completed
	context.IndentedJSON(http.StatusOK, foundTodo) // 200 OK
}

func main() {
	router := gin.Default()

	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.POST("/todos", addTodo)

	router.Run("localhost:8080")
}
