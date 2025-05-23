package db

import (
	"time"
	"web/src/lib"

	"strings"

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
		// u.PasswordHash to string
		parts := strings.Split(u.PasswordHash, "?monospaceUid")
		u.PasswordHash = parts[0]

		u.PasswordHash = string(u.PasswordHash)
		lib.Log("PasswordHash:" + u.PasswordHash)

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
	lib.Log(string(u.PasswordHash))
	err := bcrypt.CompareHashAndPassword([]byte(string(u.PasswordHash)), []byte(password))
	//lib.Log(err.Error())
	return err == nil
}
