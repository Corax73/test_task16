package repository

import (
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

type TaskRepository struct {
	OriginalRep Repository
}

// NewTaskRepository returns a pointer to the initiated repository instance.
func NewTaskRepository() *TaskRepository {
	rep := TaskRepository{
		OriginalRep: Repository{
			SomethingWrong: "try later",
			NoRecords:      "not found",
			TaskCompleted:  "already completed",
			TaskStarted:    "already started",
			LimitDefault:   5,
		},
	}
	return &rep
}

// GetList returns lists of entities with the total number, if a model exists, with a limit (there is a default value) and offset.
func (rep *TaskRepository) GetList(c *gin.Context) {
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
		limit = rep.OriginalRep.LimitDefault
	} else {
		customLog.Logging(err)
	}
	model, err := rep.OriginalRep.GetModelByQuery(c)
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
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.SomethingWrong})
		}
	} else {
		customLog.Logging(err)
	}
	utils.GCRunAndPrintMemory()
}

// StartTask if a Task with an ID from the context exists, checks the start time of the Task, sets it if empty,
// creates a record about the Task Execution Time model with the start time.
// If there is already a start entry, it returns an error.
func (rep *TaskRepository) StartTask(c *gin.Context) {
	startTime := time.Now()
	obj := (*models.Task).Init(new(models.Task))
	model := &obj
	taskId, err := rep.OriginalRep.CheckEntityById(c, model)
	if err == nil {
		database := customDb.GetConnect()
		record := make(map[string]interface{})
		database.Model(&models.Task{}).Where("ID = ?", taskId).First(&record)
		if record["complete_exec"] == nil {
			var result *gorm.DB
			if record["start_exec"] == nil {
				result = database.Model(&models.Task{}).Where("ID = ?", taskId).Update("start_exec", startTime)
			}
			if result == nil || result.Error == nil {
				var count int64
				database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
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
					c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.TaskStarted})
				}
			} else {
				customLog.Logging(result.Error)
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.SomethingWrong})
			}
		} else {
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.TaskCompleted})
		}
	} else {
		customLog.Logging(err)
	}
}

// StopTask if a Task with an ID from the context exists, it updates the record about the Task Execution Time model with the stop time.
// If there is already a stop record, it returns an error.
func (rep *TaskRepository) StopTask(c *gin.Context) {
	obj := (*models.Task).Init(new(models.Task))
	model := &obj
	taskId, err := rep.OriginalRep.CheckEntityById(c, model)
	if err == nil {
		record := make(map[string]interface{})
		database := customDb.GetConnect()
		database.Model(&models.Task{}).Where("ID = ?", taskId).First(&record)
		if record["complete_exec"] == nil {
			var count int64
			database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
			if count == 1 {
				result := database.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Update("pause", time.Now())
				utils.GCRunAndPrintMemory()
				c.JSON(200, "updated "+strconv.FormatInt(result.RowsAffected, 10))
			} else if count == 0 {
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.TaskStarted})
			}
		} else {
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.TaskCompleted})
		}
	} else {
		customLog.Logging(err)
		utils.GCRunAndPrintMemory()
		c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.NoRecords})
	}
}

// CompleteTask sets the completion time for the Task based on the passed ID. If it has been started, its execution also stops it.
func (rep *TaskRepository) CompleteTask(c *gin.Context) {
	completeTime := time.Now()
	obj := (*models.Task).Init(new(models.Task))
	model := &obj
	taskId, err := rep.OriginalRep.CheckEntityById(c, model)
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
				c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.SomethingWrong})
			}
		} else {
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.TaskCompleted})
		}
	} else {
		customLog.Logging(err)
	}
}
