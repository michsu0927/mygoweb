package db

import (
	"time"

	"gorm.io/gorm"
)

// TransactionRecord represents the transaction records table.
type TransactionRecord struct {
	gorm.Model
	UserID          string     `gorm:"column:user_id;type:char(64);not null"`
	PointsChange    int        `gorm:"not null"`
	UsedPoints      int        `gorm:"not null;default:0"`
	TransactionDate time.Time  `gorm:"column:transaction_date;type:datetime;default:CURRENT_TIMESTAMP"`
	Description     string     `gorm:"type:varchar(255)"`
	TransactionType string     `gorm:"type:varchar(255)"`
	TaskID          string     `gorm:"column:task_id;type:char(64)"`
	ExpiredDatetime *time.Time `gorm:"column:expired_datetime;type:datetime"`
}

// TableName sets the table name for TransactionRecord.
func (TransactionRecord) TableName() string {
	return "transaction_records"
}
