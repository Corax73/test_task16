package customDb

import (
	"timeTracker/customLog"
	"timeTracker/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetConnect receives data for the database from the environment file, if successful, returns the connection from the database.
func GetConnect() *gorm.DB {
	settings := utils.GetConfFromEnvFile()
	dsnStr := GetDsnString(settings)
	if dsnStr != "" {
		db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{})
		if err == nil {
			return db
		} else {
			customLog.Logging(err)
		}
	}
	return nil
}

// CloseConnect closes the connection based on the passed DB instance, returns errors when closing.
func CloseConnect(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		customLog.Logging(err)
	}
    return sqlDB.Close()
}
