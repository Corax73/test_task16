package main

import (
	"fmt"
	"time"
	"timeTracker/customDb"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID      string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	Passport  int
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "users"
}

func main() {
	database := customDb.GetConnect()
	if database != nil {
		user := User{Name: "Test", CreatedAt: time.Now(), Passport: 123}
		database.Create(&user)
		result := map[string]interface{}{}
		database.Model(&User{}).First(&result)
		fmt.Println(result)
		//	database.AutoMigrate(&User{})
	} else {
		fmt.Println(database)
	}
}
