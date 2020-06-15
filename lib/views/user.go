package views

// CUser is used to respond clients when creating User
type CUser struct {
	ID        string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UUser is used to respond clients when updating User
type UUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// RUser is used to respond clients when reading User
type RUser struct {
	ID        string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
