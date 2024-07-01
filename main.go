package main

import (
	"timeTracker/customDb"
	"timeTracker/models"

	_ "github.com/lib/pq"
)

func main() {
	database := customDb.GetConnect()
	if database != nil {
		database.AutoMigrate(&models.User{})
		database.AutoMigrate(&models.Task{})
		database.AutoMigrate(&models.TaskExecutionTime{})
		var count int64
		database.Model(&models.User{}).Count(&count)
		if count == 0 {
			customDb.SeedingUsers(database)
		}
		database.Model(&models.Task{}).Count(&count)
		if count == 0 {
			users := []map[string]interface{}{}
			database.Model(&models.User{}).Limit(10).Find(&users)
			customDb.SeedingTasks(database, users)
		}
	}
}
