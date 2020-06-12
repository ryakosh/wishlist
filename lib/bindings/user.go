package bindings

import "strings"

// CuUser is used to Create and Update User models
type CuUser struct {
	ID        string `json:"username" binding:"required,username,max=64"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=256"`
	FirstName string `json:"first_name" binding:"max=64"`
	LastName  string `json:"last_name" binding:"max=64"`
}

// RUser is used to Read User models
type RUser struct {
	ID string `json:"username" binding:"required,username,max=64"`
}

// Canonicalize UserBinding's fields to be inserted in database
func (b *CuUser) Canonicalize() {
	b.FirstName = strings.TrimSpace(b.FirstName)
	b.LastName = strings.TrimSpace(b.LastName)
}
