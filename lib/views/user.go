package views

// CUser is used to respond clients when creating/registering a User
type CUser struct {
	ID        string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UUser is used to respond clients when updating user's general information
type UUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// RUser is used to respond clients when reading general information about User
type RUser struct {
	ID        string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
