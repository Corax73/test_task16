package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primary_key, unique,type:uuid, column:id,default:uuid_generate_v4()"`
	Name           string
	CreatedAt      time.Time `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"default:NULL"`
	DeletedAt      time.Time `gorm:"default:NULL"`
	PassportNumber int
	PassportSeries int
}

// Init returns a model instance.
func (user *User) TableName() string {
	return "users"
}
