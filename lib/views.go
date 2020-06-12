package lib

type CuUserView struct {
	ID        string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RdUserView struct {
	ID        string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
