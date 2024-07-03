package router

import (
	"net/http"
	"timeTracker/customDb"
	"timeTracker/models"

	"github.com/gin-gonic/gin"
)

const noRecords = "try later"

func RunRouter() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/tasks", getTasks)
	router.Run(":4343")
}

func getUsers(c *gin.Context) {
	database := customDb.GetConnect()
	users := []map[string]interface{}{}
	database.Model(&models.User{}).Limit(10).Find(&users)
	if len(users) > 0 {
		c.JSON(200, users)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
	}
}

func getTasks(c *gin.Context) {
	database := customDb.GetConnect()
	tasks := []map[string]interface{}{}
	database.Model(&models.Task{}).Limit(10).Find(&tasks)
	if len(tasks) > 0 {
		c.JSON(200, tasks)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
	}
}
