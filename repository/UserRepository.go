package repository

import (
	"fmt"
	"net/http"
	"sort"
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
		database := customDb.GetConnect()
		result, _ := database.Debug().Migrator().ColumnTypes(&model)
		for _, v := range result {
			fieldList = append(fieldList, v.Name())
		}
		requestParams := rep.OriginalRep.GetFilterAndSortFromGetRequest(c, fieldList)
		var str strings.Builder
		if requestParams.Sorted && requestParams.SortBy != "" {
			str.WriteString(requestParams.SortBy)
			str.WriteString(" ")
			str.WriteString(requestParams.Order)
		}
		data := []map[string]interface{}{}
		database.Model(&models.Task{}).Select("id").Limit(requestParams.Limit).Offset(requestParams.Offset).Where("user_id = ?", userId).Order(str.String()).Find(&data)
		var taskIds []string
		for _, task := range data {
			id := fmt.Sprintf("%v", task["id"])
			if id != "" {
				taskIds = append(taskIds, id)
			}
		}
		if len(taskIds) > 0 {
			data := []map[string]interface{}{}
			database.Model(&models.TaskExecutionTime{}).Select("task_id", "start_exec", "pause").Limit(requestParams.Limit).
				Offset(requestParams.Offset).Where("task_id IN ? AND pause IS NOT NULL", taskIds).Order(str.String()).Find(&data)
			if len(data) > 0 {
				resp := map[string]time.Duration{}
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
								resp[task_id] += /*int(*/ dur /*)*/
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
					sortedResp := [][]interface{}{}
					keys := make([]string, 0, len(resp))
					for key := range resp {
						keys = append(keys, key)
					}

					sort.SliceStable(keys, func(i, j int) bool {
						return resp[keys[i]] > resp[keys[j]]
					})

					for _, taskId := range keys {
						duration := int(resp[taskId].Seconds())
						var hours int
						var minutes int
						var hoursStr string
						var minutesStr string
						hours = duration / 3600
						seconds := duration % 3600
						if seconds > 60 {
							minutes = seconds % 60
						}
						hoursStr = strconv.Itoa(hours)
						str.WriteString(hoursStr)
						str.WriteString(" hours")
						if minutes != 0 {
							minutesStr = strconv.Itoa(minutes)
							str.WriteString(", ")
							str.WriteString(minutesStr)
							str.WriteString(" minutes")
						}
						strDuration := str.String()
						str.Reset()
						obj := []interface{}{taskId, strDuration}
						sortedResp = append(sortedResp, obj)
					}
					utils.GCRunAndPrintMemory()
					c.JSON(200, sortedResp)
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
