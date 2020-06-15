package bindings

import "strings"

// CUser is used to Create User models
type CUser struct {
	ID        string `json:"username" binding:"required,username,max=64"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=256"`
	FirstName string `json:"first_name" binding:"max=64"`
	LastName  string `json:"last_name" binding:"max=64"`
}

// UUser is used to Update User models
type UUser struct {
	FirstName string `json:"first_name" binding:"max=64"`
	LastName  string `json:"last_name" binding:"max=64"`
}

// RUser is used to Read User models
type RUser struct {
	ID string `json:"username" binding:"required,username,max=64"`
}

// LoginUser is used for user authentication
type LoginUser struct {
	ID       string `json:"username" binding:"required,username,max=64"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

// Canonicalize UserBinding's fields to be inserted in database
func (b *CUser) Canonicalize() {
	b.FirstName = strings.TrimSpace(b.FirstName)
	b.LastName = strings.TrimSpace(b.LastName)
}
