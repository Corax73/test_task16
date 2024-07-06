package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID           uuid.UUID `gorm:"primary_key, unique,type:uuid, column:id,default:uuid_generate_v4()"`
	Title        string
	UserId       uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User         User      `gorm:"foreignKey:UserId"`
	StartExec    time.Time `gorm:"default:NULL"`
	CompleteExec time.Time `gorm:"default:NULL"`
}

func (task *Task) TableName() string {
	return "tasks"
}

// Init returns a model instance.
func (task *Task) Init() Task {
	return Task{
		ID: uuid.New(),
		Title: "stub",
		UserId: uuid.New(),
		StartExec: time.Now(),
		CompleteExec: time.Now(),
	}
}
