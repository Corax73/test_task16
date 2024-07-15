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

type GetRequestParams struct {
	Offset, Limit                      int
	Order, SortBy, FilterBy, FilterVal string
	Sorted, Filtered                   bool
}

// NewRepository returns a pointer to the initiated repository instance.
func NewRepository() *Repository {
	rep := Repository{
		SomethingWrong: "try later",
		NoRecords:      "not found",
		TaskCompleted:  "already completed",
		TaskStarted:    "already started",
		LimitDefault:   5,
	}
	return &rep
}

// GetList returns lists of entities with the total number, if a model exists, with a limit (there is a default value), offset, sort and filter by passed field.
func (rep *Repository) GetList(c *gin.Context) {
	data := []map[string]interface{}{}
	model, err := rep.GetModelByQuery(c)
	if err == nil {
		var fieldList []string
		database := customDb.GetConnect()
		result, _ := database.Debug().Migrator().ColumnTypes(&model)
		for _, v := range result {
			fieldList = append(fieldList, v.Name())
		}
		requestParams := rep.GetFilterAndSortFromGetRequest(c, fieldList)
		var count int64
		database.Model(&model).Count(&count)
		if count > 0 {
			var str strings.Builder
			if requestParams.Sorted && requestParams.SortBy != "" {
				str.WriteString(requestParams.SortBy)
				str.WriteString(" ")
				str.WriteString(requestParams.Order)
			}
			if requestParams.Filtered {
				var filterStr strings.Builder
				filterStr.WriteString(requestParams.FilterBy)
				filterStr.WriteString(" LIKE ?")
				fieldStr := filterStr.String()
				filterStr.Reset()
				filterStr.WriteString("%")
				filterStr.WriteString(requestParams.FilterVal)
				filterStr.WriteString("%")
				valStr := filterStr.String()
				database.Model(&model).Limit(requestParams.Limit).Offset(requestParams.Offset).Where(fieldStr, valStr).Order(str.String()).Find(&data)
			} else {
				database.Model(&model).Limit(requestParams.Limit).Offset(requestParams.Offset).Order(str.String()).Find(&data)
			}
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

// Delete determines the model for the route, deletes the entity using the passed ID.
func (rep *Repository) Delete(c *gin.Context) {
	model, err := rep.GetModelByQuery(c)
	if err == nil {
		id := c.Param("id")
		if id != "" {
			database := customDb.GetConnect()
			res := database.Where("id = ?", id).Delete(&model)
			if res.RowsAffected == 1 {
				utils.GCRunAndPrintMemory()
				c.JSON(200, "entity deleted")
			} else {
				customLog.Logging(res.Error)
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": rep.NoRecords})
			}
		} else {
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": rep.NoRecords})
		}
	}
}

// GetModelByQuery returns a model instance for the route from the context and an empty error, if there is no model along the route, the error will not be empty.
func (rep *Repository) GetModelByQuery(c *gin.Context) (models.Model, error) {
	var err error
	var resp models.Model
	pathSlice := strings.Split(c.Request.URL.Path, "/")
	if len(pathSlice) > 1 {
		switch pathSlice[1] {
		case "users":
			obj := (*models.User).Init(new(models.User))
			resp := &obj
			return resp, err
		case "tasks":
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
	return resp, err
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
			id, err := uuid.Parse(fmt.Sprint(paramId))
			if err == nil {
				resp = &id
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

// GetFilterAndSortFromGetRequest gets the limit, offset, filtering and sorting parameters from the request,
// checks the matches of the names of the sorting and filtering fields with the passed slice of table names of the model fields,
// returns a structure with values.
func (rep *Repository) GetFilterAndSortFromGetRequest(c *gin.Context, fieldList []string) GetRequestParams {
	var offset int
	var limit int
	sortBy := "id"
	order := "desc"
	var sorted bool
	var filtered bool
	var filterBy string
	var filterVal string
	var resp GetRequestParams

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

	sort := c.Query("sort")
	if sort != "" {
		splits := strings.Split(sort, ".")
		requestField, requestOrder := splits[0], splits[1]
		if requestOrder != "desc" && requestOrder != "asc" {
			order = "desc"
		} else {
			order = requestOrder
		}
		if slices.Contains(fieldList, requestField) {
			sorted = true
			sortBy = requestField
		}
	}

	filter := c.Query("filter")
	if filter != "" {
		splits := strings.Split(filter, ".")
		requestField, requestValue := splits[0], splits[1]
		if slices.Contains(fieldList, requestField) && requestValue != "" {
			filtered = true
			filterBy = requestField
			filterVal = requestValue
		}
	}

	resp = GetRequestParams{
		Offset:    offset,
		Limit:     limit,
		Order:     order,
		SortBy:    sortBy,
		FilterBy:  filterBy,
		FilterVal: filterVal,
		Sorted:    sorted,
		Filtered:  filtered,
	}
	return resp
}
