package customDb

import (
	"strconv"
	"time"
	"timeTracker/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Creates test entries in the users table.
func SeedingUsers(db *gorm.DB) {
	for i := 0; i < 10; i++ {
		id := uuid.New()
		user := models.User{ID: id, Name: "Test" + strconv.Itoa(i), CreatedAt: time.Now(), PassportNumber: 123 + i, PassportSeries: 321123 + i}
		db.Create(&user)
	}
}
