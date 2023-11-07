package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string  `gorm:"index:email,unique;not null" json:"email"`
	FirstName     string  `gorm:"column:first_name;not null" json:"first_name"`
	SecondName    string  `gorm:"column:second_name;not null" json:"second_name"`
	Username      string  `gorm:"index:username,unique;not null" json:"username"`
	Password      string  `gorm:"column:password; not null" json:"password"`
	Phone         *string `gorm:"index:phone;unique" json:"phone"`
	FirebaseToken string  `gorm:"column:firebase_token" json:"firebase_token"`
}

type RegisterRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
