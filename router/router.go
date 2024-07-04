package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"timeTracker/customDb"
	"timeTracker/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const somethingWrong = "try later"
const noRecords = "not found"

func RunRouter() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/tasks", getTasks)
	router.POST("/tasks/start", startTask)
	router.POST("/tasks/stop", stopTask)
	router.Run(":4343")
}

func getUsers(c *gin.Context) {
	database := customDb.GetConnect()
	users := []map[string]interface{}{}
	database.Model(&models.User{}).Limit(10).Find(&users)
	if len(users) > 0 {
		c.JSON(200, users)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": somethingWrong})
	}
}

func getTasks(c *gin.Context) {
	database := customDb.GetConnect()
	tasks := []map[string]interface{}{}
	database.Model(&models.Task{}).Limit(10).Find(&tasks)
	if len(tasks) > 0 {
		c.JSON(200, tasks)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": somethingWrong})
	}
}

func startTask(c *gin.Context) {
	paramId := c.DefaultPostForm("id", "0")
	if paramId == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
	} else {
		database := customDb.GetConnect()
		var count int64
		database.Model(&models.Task{}).Where("id = ?", paramId).Count(&count)
		if count > 0 {
			taskId, err := uuid.Parse(fmt.Sprint(paramId))
			if err == nil {
				database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
				if count == 0 {
					database.Save(&models.TaskExecutionTime{ID: uuid.New(), TaskId: taskId, StartExec: time.Now()})
					c.JSON(200, "started")
				} else if count > 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "already started"})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
		}
	}
}

func stopTask(c *gin.Context) {
	paramId := c.DefaultPostForm("id", "0")
	if paramId == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
	} else {
		database := customDb.GetConnect()
		var count int64
		database.Model(&models.Task{}).Where("id = ?", paramId).Count(&count)
		if count > 0 {
			taskId, err := uuid.Parse(fmt.Sprint(paramId))
			if err == nil {
				database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
				if count == 1 {
					result := database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Update("pause", time.Now())
					c.JSON(200, "updated "+strconv.FormatInt(result.RowsAffected, 10))
				} else if count == 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "already stopped"})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
		}
	}
}
