package repository

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
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

// GetList returns lists of entities with the total number, if a model exists, with a limit (there is a default value), offset and sort by passed field.
func (rep *Repository) GetList(c *gin.Context) {
	database := customDb.GetConnect()
	data := []map[string]interface{}{}
	var offset int
	var limit int
	requestOffset := c.Query("offset")
	if requestOffset != "" {
		val, err := strconv.Atoi(requestOffset)
		if err != nil {
			offset = 0
			customLog.Logging(err)
		} else {
			offset = val
		}
	}
	requestLimit := c.Query("limit")
	if requestLimit != "" {
		val, err := strconv.Atoi(requestLimit)
		if err != nil {
			limit = rep.LimitDefault
			customLog.Logging(err)
		} else {
			limit = val
		}
	} else {
		limit = rep.LimitDefault
	}
	model, err := rep.GetModelByQuery(c)
	if err == nil {
		var fieldList []string
		field := "id"
		order := "desc"
		var sorted bool
		result, _ := database.Debug().Migrator().ColumnTypes(&model)
		for _, v := range result {
			fieldList = append(fieldList, v.Name())
		}
		sort := c.Query("sort")
		if sort != "" {
			splits := strings.Split(sort, ".")
			requestField, requestOrder := splits[0], splits[1]
			if requestOrder != "desc" && requestOrder != "asc" {
				order = "desc"
			} else {
				order = requestOrder
			}
			if slices.Contains(fieldList, field) {
				sorted = true
				field = requestField
			}
		}
		var count int64
		database.Model(&model).Count(&count)
		if count > 0 {
			var str strings.Builder
			if sorted && field != "" {
				str.WriteString(field)
				str.WriteString(" ")
				str.WriteString(order)
			}
			database.Model(&model).Limit(limit).Offset(offset).Order(str.String()).Find(&data)
			total := make(map[string]interface{})
			total["total"] = count
			data = append(data, total)
			utils.GCRunAndPrintMemory()
			c.JSON(200, data)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.SomethingWrong})
		}
	} else {
		customLog.Logging(err)
	}
	utils.GCRunAndPrintMemory()
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
