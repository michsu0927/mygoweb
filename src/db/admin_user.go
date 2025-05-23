package db

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminUser model
type AdminUser struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BeforeSave hook to hash password before saving
func (u *AdminUser) BeforeSave(tx *gorm.DB) (err error) {
	if u.PasswordHash != "" { // Only hash if password is set/changed
		//print PasswordHash
		//fmt.Println("PasswordHash: ", u.PasswordHash)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashedPassword)
	}
	return
}

// CheckPassword checks if the provided password is correct
func (u *AdminUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
