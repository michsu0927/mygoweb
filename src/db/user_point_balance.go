package db

import "gorm.io/gorm"

// UserPointBalance represents the user point balance table.
type UserPointBalance struct {
	gorm.Model
	UserID  string `gorm:"column:user_id;type:char(64);not null;uniqueIndex"`
	Balance int64  `gorm:"column:balance;not null;default:0"`
}

// TableName sets the table name for UserPointBalance.
func (UserPointBalance) TableName() string {
	return "user_point_balance"
}
