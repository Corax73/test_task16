package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"timeTracker/customDb"
	"timeTracker/customLog"
	"timeTracker/models"
	"timeTracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const somethingWrong string = "try later"
const noRecords string = "not found"
const limitDefault int = 5

func RunRouter() {
	utils.GCRunAndPrintMemory()
	router := gin.Default()
	router.GET("/users", getList)
	router.GET("/tasks", getList)
	router.POST("/tasks/start", startTask)
	router.POST("/tasks/stop", stopTask)
	router.POST("/tasks/complete", completeTask)
	router.Run(":4343")
}

// getList returns lists of entities with the total number, if a model exists, with a limit (there is a default value) and offset.
func getList(c *gin.Context) {
	database := customDb.GetConnect()
	data := []map[string]interface{}{}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	} else {
		customLog.Logging(err)
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = limitDefault
	} else {
		customLog.Logging(err)
	}
	model, err := getModelByQuery(c)
	if err == nil {
		var count int64
		database.Model(&model).Count(&count)
		if count > 0 {
			database.Model(&model).Limit(limit).Offset(offset).Find(&data)
			total := make(map[string]interface{})
			total["total"] = count
			data = append(data, total)
			utils.GCRunAndPrintMemory()
			c.JSON(200, data)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": somethingWrong})
		}
	} else {
		customLog.Logging(err)
	}
	utils.GCRunAndPrintMemory()
}

// getModelByQuery returns a model instance for the route from the context and an empty error, if there is no model along the route, the error will not be empty.
func getModelByQuery(c *gin.Context) (models.Model, error) {
	var err error
	switch c.Request.URL.Path {
	case "/users":
		obj := (*models.User).Init(new(models.User))
		resp := &obj
		return resp, err
	case "/tasks":
		obj := (*models.Task).Init(new(models.Task))
		resp := &obj
		return resp, err
	default:
		obj := (*models.User).Init(new(models.User))
		resp := &obj
		err = errors.New("unknown route")
		customLog.Logging(err)
		return resp, err
	}
}

// startTask if a Task with an ID from the context exists, checks the start time of the Task, sets it if empty,
// creates a record about the Task Execution Time model with the start time.
// If there is already a start entry, it returns an error.
func startTask(c *gin.Context) {
	startTime := time.Now()
	obj := (*models.Task).Init(new(models.Task))
	model := &obj
	taskId, err := checkEntityById(c, model)
	if err == nil {
		var count int64
		database := customDb.GetConnect()
		database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
		record := make(map[string]interface{})
		database.Model(&models.Task{}).Where("ID = ?", taskId).First(&record)
		var result *gorm.DB
		if record["start_exec"] == nil {
			result = database.Model(&models.Task{}).Where("ID = ?", taskId).Update("start_exec", startTime)
		}
		if result == nil || result.Error == nil {
			if count == 0 {
				taskId, err := uuid.Parse(fmt.Sprint(taskId))
				if err == nil {
					database.Save(&models.TaskExecutionTime{ID: uuid.New(), TaskId: taskId, StartExec: startTime})
					utils.GCRunAndPrintMemory()
					c.JSON(200, "started")
				} else {
					customLog.Logging(err)
					utils.GCRunAndPrintMemory()
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
			} else if count > 0 {
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": "already started"})
			}
		} else {
			customLog.Logging(result.Error)
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": somethingWrong})
		}
	} else {
		customLog.Logging(err)
	}
}

// stopTask if a Task with an ID from the context exists, it updates the record about the Task Execution Time model with the stop time.
// If there is already a stop record, it returns an error.
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
			utils.GCRunAndPrintMemory()
			c.JSON(200, "updated "+strconv.FormatInt(result.RowsAffected, 10))
		} else if count == 0 {
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": "already stopped"})
		}
	} else {
		customLog.Logging(err)
		utils.GCRunAndPrintMemory()
		c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
	}
}

// completeTask sets the completion time for the Task based on the passed ID. If it has been started, its execution also stops it.
func completeTask(c *gin.Context) {
	completeTime := time.Now()
	obj := (*models.Task).Init(new(models.Task))
	model := &obj
	taskId, err := checkEntityById(c, model)
	if err == nil {
		database := customDb.GetConnect()
		record := make(map[string]interface{})
		database.Model(&models.Task{}).Where("ID = ?", taskId).First(&record)
		var result *gorm.DB
		if record["complete_exec"] == nil {
			result = database.Model(&models.Task{}).Where("ID = ?", taskId).Update("complete_exec", completeTime)
			if result.Error == nil {
				var count int64
				database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
				if count == 1 {
					result := database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Update("pause", completeTime)
					utils.GCRunAndPrintMemory()
					c.JSON(200, "updated "+strconv.FormatInt(result.RowsAffected, 10))
				} else if count == 0 {
					c.JSON(200, "updated "+strconv.FormatInt(result.RowsAffected, 10))
					utils.GCRunAndPrintMemory()
				}
			} else {
				customLog.Logging(result.Error)
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": somethingWrong})
			}
		}
	} else {
		customLog.Logging(err)
	}
}

// checkEntityById using the ID from the passed context, searches for a record based on the passed model. If exists, returns the model *uuid.UUID and an empty error.
// Otherwise default *uuid.UUID and non-empty error.
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
				customLog.Logging(err)
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			err = errors.New(noRecords + " " + model.TableName())
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": noRecords})
		}
	}
	return resp, err
}
