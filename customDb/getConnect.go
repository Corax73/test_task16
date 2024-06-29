package customDb

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Receives data for the database from the environment file, if successful, returns the connection from the database.
func GetConnect() *gorm.DB {
	dsnMap := GetConfFromEnvFile()
	dsnStr := GetDsnString(dsnMap)
	if dsnStr != "" {
		db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{})
		if err == nil {
			return db
		} else {
			fmt.Println(err)
		}
	}
	return nil
}
