package customDb

import (
	"timeTracker/models"
	"timeTracker/utils"
)

// Init conducts initial migrations and populates test data. Returns true on success.
func Init() bool {
	var resp bool
	database := GetConnect()
	if database != nil {
		errUser := database.AutoMigrate(&models.User{})
		errTask := database.AutoMigrate(&models.Task{})
		errTaskExecutionTime := database.AutoMigrate(&models.TaskExecutionTime{})
		if errUser == nil && errTask == nil && errTaskExecutionTime == nil {
			var count int64
			database.Model(&models.User{}).Count(&count)
			if count == 0 {
				SeedingUsers(database)
			}
			database.Model(&models.Task{}).Count(&count)
			if count == 0 {
				users := []map[string]interface{}{}
				database.Model(&models.User{}).Limit(10).Find(&users)
				SeedingTasks(database, users)
				database.Model(&models.Task{}).Count(&count)
			}
			if count > 0 {
				resp = true
			}
		}
	}
	utils.PrintMemoryAndGC()
	return resp
}
