package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
}

// SetPassword 设置密码 (加密)
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type VerificationCode struct {
	gorm.Model
	Email     string    `json:"email" gorm:"not null"`
	Code      string    `json:"code" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}
