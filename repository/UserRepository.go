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
	Repository
}

// NewTaskRepository returns a pointer to the initiated repository instance.
func NewUserRepository() *UserRepository {
	rep := UserRepository{
		Repository{
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
			tx := database.Begin()
			result := tx.Create(&user)
			if result.Error == nil {
				res := tx.Commit()
				if res.Error == nil {
					utils.GCRunAndPrintMemory()
					c.JSON(200, "created")
				} else {
					tx.Rollback()
					customLog.Logging(res.Error)
					utils.GCRunAndPrintMemory()
					c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
				}
			} else {
				tx.Rollback()
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
	userId, err := rep.CheckEntityById(c, model)
	if err == nil {
		database := customDb.GetConnect()
		data := []map[string]interface{}{}
		database.Model(&models.Task{}).Select("id").Where("user_id = ?", userId).Find(&data)
		var taskIds []string
		for _, task := range data {
			id := fmt.Sprintf("%v", task["id"])
			if id != "" {
				taskIds = append(taskIds, id)
			}
		}
		if len(taskIds) > 0 {
			data := []map[string]interface{}{}
			database.Model(&models.TaskExecutionTime{}).Select("task_id", "start_exec", "pause").Where("task_id IN ? AND pause IS NOT NULL", taskIds).Find(&data)
			if len(data) > 0 {
				resp := map[string]time.Duration{}
				for _, task := range data {
					task_id := fmt.Sprintf("%v", task["task_id"])
					if task_id != "" {
						const layout = "2006-01-02 15:04:05"
						startExec := fmt.Sprintf("%v", task["start_exec"])
						startExec = startExec[:len(startExec)-10]
						parseTimeStart, err := time.Parse(layout, startExec)
						fmt.Println("parseTimeStart=", parseTimeStart)
						if err == nil {
							pause := fmt.Sprintf("%v", task["pause"])
							pause = pause[:len(pause)-10]
							parseTimePause, err := time.Parse(layout, pause)
							fmt.Println("parseTimePause=", parseTimePause)
							if err == nil {
								dur := parseTimePause.Sub(parseTimeStart)
								fmt.Println("dur=", dur)
								resp[task_id] += dur
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
				fmt.Println("resp=", resp)
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
							minutes = seconds / 60
						}
						fmt.Println("hours=", hours)
						fmt.Println("minutes=", minutes)
						fmt.Println("seconds=", seconds)
						hoursStr = strconv.Itoa(hours)
						var str strings.Builder
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
