package model

import "github.com/ryakosh/wishlist/lib/db"

type User struct {
	ID             string  `json:"id"`
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Wishes         string  `json:"wishes"`
	Friends        string  `json:"friends"`
	FriendRequests string  `json:"friendRequests"`
}

type Users struct {
	Query         string      `json:"users"` // User ID as string or wish ID as string
	Count         int         `json:"count"`
	InObj         interface{} // Parent model
	InAssociation db.Association
}

type NewUser struct {
	ID        string  `json:"id" validate:"username,max=64"`
	FirstName *string `json:"firstName" validate:"omitempty,max=64"`
	LastName  *string `json:"lastName" validate:"omitempty,max=64"`
	Email     string  `json:"email" validate:"email"`
	Password  string  `json:"password" validate:"min=8,max=256"`
}

type UpdateUser struct {
	FirstName *string `json:"firstName" validate:"omitempty,max=64"`
	LastName  *string `json:"lastName" validate:"omitempty,max=64"`
}

type Login struct {
	ID       string `json:"id" validate:"username,max=64"`
	Password string `json:"password" validate:"min=8,max=256"`
}
