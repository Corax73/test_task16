package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"timeTracker/customDb"
	"timeTracker/models"
	"timeTracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const somethingWrong = "try later"
const noRecords = "not found"

func RunRouter() {
	utils.PrintMemoryAndGC()
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
	utils.PrintMemoryAndGC()
}

func getTasks(c *gin.Context) {
	database := customDb.GetConnect()
	tasks := []map[string]interface{}{}
	database.Model(&models.Task{}).Limit(10).Find(&tasks)
	if len(tasks) > 0 {
		utils.PrintMemoryAndGC()
		c.JSON(200, tasks)
	} else {
		utils.PrintMemoryAndGC()
		c.JSON(http.StatusBadRequest, gin.H{"error": somethingWrong})
	}
}

func startTask(c *gin.Context) {
	obj := (*models.Task).Init(new(models.Task))
	model := &obj
	taskId, err := checkEntityById(c, model)
	if err == nil {
		var count int64
		database := customDb.GetConnect()
		database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
		if count == 0 {
			taskId, err := uuid.Parse(fmt.Sprint(taskId))
			if err == nil {
				database.Save(&models.TaskExecutionTime{ID: uuid.New(), TaskId: taskId, StartExec: time.Now()})
				utils.PrintMemoryAndGC()
				c.JSON(200, "started")
			} else {
				utils.PrintMemoryAndGC()
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else if count > 0 {
			utils.PrintMemoryAndGC()
			c.JSON(http.StatusBadRequest, gin.H{"error": "already started"})
		}
	}
}

func stopTask(c *gin.Context) {
	obj := (*models.Task).Init(new(models.Task))
	model := &obj
	taskId, err := checkEntityById(c, model)
	if err == nil {
		var count int64
		database := customDb.GetConnect()
		database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
		if count == 1 {
			result := database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Update("pause", time.Now())
			utils.PrintMemoryAndGC()
			c.JSON(200, "updated "+strconv.FormatInt(result.RowsAffected, 10))
		} else if count == 0 {
			utils.PrintMemoryAndGC()
			c.JSON(http.StatusBadRequest, gin.H{"error": "already stopped"})
		}
	} else {
		utils.PrintMemoryAndGC()
		c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
	}
}

func checkEntityById(c *gin.Context, model models.Model) (*uuid.UUID, error) {
	var err error
	defaultId := uuid.New()
	resp := &defaultId
	paramId := c.DefaultPostForm("id", "0")
	if paramId == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paramId = 0"})
	} else {
		database := customDb.GetConnect()
		var count int64
		database.Model(&model).Where("id = ?", paramId).Count(&count)
		if count > 0 {
			taskId, err := uuid.Parse(fmt.Sprint(paramId))
			if err == nil {
				resp = &taskId
			} else {
				utils.PrintMemoryAndGC()
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			utils.PrintMemoryAndGC()
			c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
		}
	}
	return resp, err
}
