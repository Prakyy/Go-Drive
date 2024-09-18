package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email		string `gorm:"unique;not null;type:varchar(100);default:null"`
	Password	string `gorm:"not null;default:null"`
}