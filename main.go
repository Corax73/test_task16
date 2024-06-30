package main

import (
	"fmt"
	"timeTracker/customDb"
	"timeTracker/models"

	_ "github.com/lib/pq"
)

func main() {
	database := customDb.GetConnect()
	if database != nil {
		database.AutoMigrate(&models.User{})
		var count int64
		database.Model(&models.User{}).Count(&count)
		if count == 0 {
			customDb.SeedingUsers(database)
		}
		results := []map[string]interface{}{}
		database.Model(&models.User{}).Limit(10).Find(&results)
		fmt.Print(len(results))
	}
}
