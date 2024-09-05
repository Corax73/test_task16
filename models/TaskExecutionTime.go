package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskExecutionTime struct {
	Task      Task      `gorm:"foreignKey:TaskId"`
	StartExec time.Time `gorm:"type:TIMESTAMP;null;default:null"`
	Pause     time.Time `gorm:"type:TIMESTAMP;null;default:null"`
	ID        uuid.UUID `gorm:"primary_key, unique,type:uuid, column:id,default:uuid_generate_v4()"`
	TaskId    uuid.UUID `json:"task_id" gorm:"type:uuid"`
}

func (task *TaskExecutionTime) TableName() string {
	return "tasks_exec_time"
}

// Init returns a model instance.
func (task *TaskExecutionTime) Init() TaskExecutionTime {
	return TaskExecutionTime{
		ID:        uuid.New(),
		TaskId:    uuid.New(),
		StartExec: time.Now(),
		Pause:     time.Now(),
	}
}
