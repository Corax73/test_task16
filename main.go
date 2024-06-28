package main

import (
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
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
	database := getConnect()
	if database != nil {
		database.AutoMigrate(&User{})
	} else {
		fmt.Println(database)
	}
}

func getConnect() *gorm.DB {
	db, err := gorm.Open(postgres.Open("user=postgres password=postgres dbname=postgres sslmode=disable"), &gorm.Config{})
	if err == nil {
		return db
	} else {
		fmt.Println(err)
	}
	return nil
}
