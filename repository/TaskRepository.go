package repository

import (
	"fmt"
	"net/http"
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
			TaskStopped:    "already stopped",
			LimitDefault:   5,
		},
	}
	return &rep
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
						tx := database.Begin()
						res := tx.Save(&models.TaskExecutionTime{ID: uuid.New(), TaskId: taskId, StartExec: startTime})
						if res.Error == nil {
							res := tx.Commit()
							if res.Error == nil {
								utils.GCRunAndPrintMemory()
								c.JSON(200, "started")
							} else {
								tx.Rollback()
								customLog.Logging(res.Error)
								utils.GCRunAndPrintMemory()
								c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
							}
						} else {
							tx.Rollback()
							customLog.Logging(res.Error)
							utils.GCRunAndPrintMemory()
							c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
						}
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
				tx := database.Begin()
				res := tx.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Update("pause", time.Now())
				if res.Error == nil {
					res := tx.Commit()
					if res.Error == nil {
						utils.GCRunAndPrintMemory()
						c.JSON(200, "task stopped")
					} else {
						tx.Rollback()
						customLog.Logging(res.Error)
						utils.GCRunAndPrintMemory()
						c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
					}
				} else {
					tx.Rollback()
					customLog.Logging(res.Error)
					utils.GCRunAndPrintMemory()
					c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
				}
			} else if count == 0 {
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.TaskStopped})
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
		if record["complete_exec"] == nil {
			tx := database.Begin()
			res := tx.Model(&models.Task{}).Where("ID = ?", taskId).Update("complete_exec", completeTime)
			if res.Error == nil {
				if res.Error == nil {
					var count int64
					tx.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Count(&count)
					if count == 1 {
						res := tx.Model(&models.TaskExecutionTime{}).Where("task_id = ? AND pause IS NULL", taskId).Update("pause", completeTime)
						if res.Error == nil {
							res := tx.Commit()
							if res.Error == nil {
								utils.GCRunAndPrintMemory()
								c.JSON(200, "task completed")
							} else {
								tx.Rollback()
								customLog.Logging(res.Error)
								utils.GCRunAndPrintMemory()
								c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
							}
						} else {
							tx.Rollback()
							customLog.Logging(res.Error)
							utils.GCRunAndPrintMemory()
							c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
						}
					} else if count == 0 {
						c.JSON(200, "task completed")
						utils.GCRunAndPrintMemory()
					}
				} else {
					tx.Rollback()
					customLog.Logging(res.Error)
					utils.GCRunAndPrintMemory()
					c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
				}
			} else {
				tx.Rollback()
				customLog.Logging(res.Error)
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
			}
		} else {
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.OriginalRep.TaskCompleted})
		}
	} else {
		customLog.Logging(err)
	}
}
