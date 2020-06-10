package main

import (
	"time"
)

type User struct {
	ID                string     `json:"username" gorm:"type:varchar(64)" binding:"required,username,max=64"`
	Email             string     `json:"email" gorm:"varchar(254)" binding:"required,email"`
	IsEmailVerified   bool       `json:"is_email_verfied" binding:"-"`
	Password          string     `json:"password" gorm:"varchar(256)" binding:"required,min=8,max=256"`
	FirstName         string     `json:"first_name" gorm:"type:varchar(64)" binding:"max=64"`
	LastName          string     `json:"last_name" gorm:"type:varchar(64)" binding:"max=64"`
	Wishes            []Wish     `json:"wishes" binding:"-"`
	FulfilledWishes   []Wish     `json:"fulfilled_wishes" gorm:"foreignkey:FulfilledBy" binding:"-"`
	WantFulfillWishes []Wish     `json:"want_fulfill_wishes" gorm:"many2many:userswant_wishes" binding:"-"`
	CreatedAt         *time.Time `binding:"-"`
	UpdatedAt         *time.Time `binding:"-"`
}

type Wish struct {
	ID          uint       `json:"id" binding:"-"`
	UserID      string     `json:"username" binding:"-"`
	Name        string     `json:"name" gorm:"type:varchar(256)" binding:"required,max=256"`
	Description string     `json:"description" gorm:"type:text"`
	Link        string     `json:"link" binding:"url"`
	Image       string     `json:"image" binding:"url"`
	FulfilledBy string     `json:"fulfilled_by" binding:"-"`
	CreatedAt   *time.Time `binding:"-"`
	UpdatedAt   *time.Time `binding:"-"`
}
