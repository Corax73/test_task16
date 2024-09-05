package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	User         User      `gorm:"foreignKey:UserId"`
	StartExec    time.Time `gorm:"type:TIMESTAMP;null;default:null"`
	CompleteExec time.Time `gorm:"type:TIMESTAMP;null;default:null"`
	ID           uuid.UUID `gorm:"primary_key, unique,type:uuid, column:id,default:uuid_generate_v4()"`
	UserId       uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Title        string
}

func (task *Task) TableName() string {
	return "tasks"
}

// Init returns a model instance.
func (task *Task) Init() Task {
	return Task{
		ID:           uuid.New(),
		Title:        "stub",
		UserId:       uuid.New(),
		StartExec:    time.Now(),
		CompleteExec: time.Now(),
	}
}
