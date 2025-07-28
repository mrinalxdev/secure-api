package models 

import "gorm.io/gorm"

type User struct {
	gorm.Model 

	Username string `json:"username" gorm:"uniqueIndex;not null`
	Password string `json:"-" gorm:"not null"`
	Role string `json:"role" gorm:"default:'user'"`

}