package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskExecutionTime struct {
	ID        uuid.UUID `gorm:"primary_key, unique,type:uuid, column:id,default:uuid_generate_v4()"`
	TaskId    uuid.UUID `json:"task_id" gorm:"type:uuid"`
	Task      Task      `gorm:"foreignKey:TaskId"`
	StartExec time.Time `gorm:"default:NULL"`
	Pause     time.Time `gorm:"default:NULL"`
}

func (task *TaskExecutionTime) TableName() string {
	return "tasks_exec_time"
}
