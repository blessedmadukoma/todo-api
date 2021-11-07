package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string
	Lastname  string
	Email     string `gorm:"uniqueIndex"` // unique key
	Password  string
	Todos     []Todo
}

type Todo struct {
	gorm.Model
	Text   string
	UserID uint
}

type Response struct {
	Status  string
	Message string
	// Data    interface{}
	Data map[string]interface{}
}
