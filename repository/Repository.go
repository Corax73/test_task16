package repository

import (
	"errors"
	"fmt"
	"net/http"
	"timeTracker/customDb"
	"timeTracker/customLog"
	"timeTracker/models"
	"timeTracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Repository struct {
	SomethingWrong, NoRecords, TaskCompleted, TaskStarted string
	LimitDefault                                          int
}

// GetModelByQuery returns a model instance for the route from the context and an empty error, if there is no model along the route, the error will not be empty.
func (rep *Repository) GetModelByQuery(c *gin.Context) (models.Model, error) {
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

// CheckEntityById using the ID from the passed context, searches for a record based on the passed model. If exists, returns the model *uuid.UUID and an empty error.
// Otherwise default *uuid.UUID and non-empty error.
func (rep *Repository) CheckEntityById(c *gin.Context, model models.Model) (*uuid.UUID, error) {
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
			err = errors.New(rep.NoRecords + " " + model.TableName())
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.NoRecords})
		}
	}
	return resp, err
}
