package db

import (
	"time"

	"gorm.io/gorm"
)

// Task represents the task table.
type Task struct {
	gorm.Model
	TaskID          string     `gorm:"column:task_id;type:char(64);primaryKey"`
	UserID          string     `gorm:"column:user_id;type:char(64)"`
	TaskType        string     `gorm:"type:varchar(255)"`
	Description     string     `gorm:"type:varchar(255)"`
	PointsChange    int        `gorm:"not null"`
	Status          int        `gorm:"type:tyint;not null;default:0"`
	ExpiredDatetime *time.Time `gorm:"column:expired_datetime;type:datetime"`
	//CreateDatetime  time.Time  `gorm:"column:create_datetime;type:datetime;default:CURRENT_TIMESTAMP"`
}

// TableName sets the table name for Task.
func (Task) TableName() string {
	return "task"
}
