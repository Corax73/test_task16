package customDb

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"timeTracker/customLog"
	"timeTracker/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedingTasks creates test entries in the task table based on the passed slice of user maps.
func SeedingTasks(db *gorm.DB, userData []map[string]interface{}) {
	for i := 0; i < 10; i++ {
		id := uuid.New()
		completed := rand.Intn(2)
		startExec := time.Now()
		var completeExec time.Time
		if completed == 1 {
			completeExec = startExec.Add(time.Hour*1 + time.Minute*15)
		}
		userId, err := uuid.Parse(fmt.Sprint(userData[i]["id"]))
		if err != nil {
			fmt.Println(userId)
			task := models.Task{ID: id, Title: "Test_task" + strconv.Itoa(i), StartExec: startExec, CompleteExec: completeExec, UserId: userId}
			db.Create(&task)
		} else {
			customLog.Logging(err)
		}
	}
}
