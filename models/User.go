package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primary_key, unique,type:uuid, column:id,default:uuid_generate_v4()"`
	Name           string    `gorm:"unique"`
	CreatedAt      time.Time `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"default:NULL"`
	DeletedAt      time.Time `gorm:"default:NULL"`
	PassportNumber int
	PassportSeries int
}

func (user *User) TableName() string {
	return "users"
}

// Init returns a model instance.
func (task *User) Init() User {
	return User{
		ID:             uuid.New(),
		Name:           "Default name",
		CreatedAt:      time.Now(),
		PassportNumber: 123,
		PassportSeries: 321456,
	}
}
