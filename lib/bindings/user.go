package bindings

// CUser is used to create/register a new user
type CUser struct {
	ID        string `json:"username" binding:"required,username,max=64"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=256"`
	FirstName string `json:"first_name" binding:"max=64"`
	LastName  string `json:"last_name" binding:"max=64"`
}

// UUser is used to update user's general information
type UUser struct {
	FirstName string `json:"first_name" binding:"max=64"`
	LastName  string `json:"last_name" binding:"max=64"`
}

// RUser is used to read user's general information
type RUser struct {
	ID string `json:"username" binding:"required,username,max=64"`
}

// LoginUser is used for user authentication
type LoginUser struct {
	ID       string `json:"username" binding:"required,username,max=64"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}
