package views

// CuUser is used to respond clients when creating and updating User
type CuUser struct {
	ID        string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// RUser is used to respond clients when reading User
type RUser struct {
	ID        string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
