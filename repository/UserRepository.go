package repository

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"timeTracker/customDb"
	"timeTracker/customLog"
	"timeTracker/models"
	"timeTracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRepository struct {
	OriginalRep Repository
}

// NewTaskRepository returns a pointer to the initiated repository instance.
func NewUserRepository() *UserRepository {
	rep := UserRepository{
		OriginalRep: Repository{
			SomethingWrong: "try later",
			NoRecords:      "not found",
			LimitDefault:   5,
		},
	}
	return &rep
}

// Create based on the data from the request, if they are not empty, creates a User record and returns it and an empty error,
// otherwise returns an empty User instance and a non-empty error.
func (rep *UserRepository) Create(c *gin.Context) (*models.User, error) {
	var err error
	user := models.User{}
	resp := &user
	newId := uuid.New()
	name := c.DefaultPostForm("name", "")
	passportNumber, errNumber := strconv.Atoi(c.DefaultPostForm("passportNumber", "0"))
	passportSeries, errSeries := strconv.Atoi(c.DefaultPostForm("passportSeries", "0"))
	if errNumber == nil && errSeries == nil {
		if name != "" && passportNumber != 0 && passportSeries != 0 {
			database := customDb.GetConnect()
			user := models.User{ID: newId, Name: name, PassportNumber: passportNumber, PassportSeries: passportSeries}
			result := database.Create(&user)
			if result.Error == nil {
				utils.GCRunAndPrintMemory()
				c.JSON(200, "created")
			} else {
				customLog.Logging(result.Error)
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			}
			resp = &user
		}
	} else {
		if errNumber != nil {
			err = errNumber
		} else {
			err = errSeries
		}
	}
	return resp, err
}

func (rep *UserRepository) GetTaskExecutionTime(c *gin.Context) {
	obj := (*models.User).Init(new(models.User))
	model := &obj
	userId, err := rep.OriginalRep.CheckEntityById(c, model)
	if err == nil {
		var fieldList []string
		var offset int
		var limit int
		var sorted bool
		sortBy := "id"
		order := "desc"
		database := customDb.GetConnect()
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
				limit = rep.OriginalRep.LimitDefault
				customLog.Logging(err)
			} else {
				limit = val
			}
		} else {
			limit = rep.OriginalRep.LimitDefault
		}
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
			if slices.Contains(fieldList, requestField) {
				sorted = true
				sortBy = requestField
			}
		}
		var str strings.Builder
		if sorted && sortBy != "" {
			str.WriteString(sortBy)
			str.WriteString(" ")
			str.WriteString(order)
		}
		data := []map[string]interface{}{}
		database.Model(&models.Task{}).Select("id").Limit(limit).Offset(offset).Where("user_id = ?", userId).Order(str.String()).Find(&data)
		var taskIds []string
		for _, task := range data {
			id := fmt.Sprintf("%v", task["id"])
			if id != "" {
				taskIds = append(taskIds, id)
			}
		}
		if len(taskIds) > 0 {
			data := []map[string]interface{}{}
			database.Model(&models.TaskExecutionTime{}).Select("task_id", "start_exec", "pause").Limit(limit).
				Offset(offset).Where("task_id IN ? AND pause IS NOT NULL", taskIds).Order(str.String()).Find(&data)
			if len(data) > 0 {
				resp := map[string]int{}
				for _, task := range data {
					task_id := fmt.Sprintf("%v", task["task_id"])
					if task_id != "" {
						const layout = "2006-01-02 15:04:05"
						startExec := fmt.Sprintf("%v", task["start_exec"])
						startExec = startExec[:len(startExec)-10]
						parseTimeStart, err := time.Parse(layout, startExec)
						if err == nil {
							pause := fmt.Sprintf("%v", task["pause"])
							pause = pause[:len(pause)-10]
							parseTimePause, err := time.Parse(layout, pause)
							if err == nil {
								dur := parseTimePause.Sub(parseTimeStart)
								resp[task_id] += int(dur)
							} else {
								utils.GCRunAndPrintMemory()
								customLog.Logging(err)
							}
						} else {
							utils.GCRunAndPrintMemory()
							customLog.Logging(err)
						}
					}
				}
				if len(resp) > 0 {
					utils.GCRunAndPrintMemory()
					c.JSON(200, resp)
				}
			} else {
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": "no tasks found for user"})
			}
		} else {
			utils.GCRunAndPrintMemory()
			c.JSON(http.StatusBadRequest, gin.H{"error": "no tasks found for user"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
	}
}
