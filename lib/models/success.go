package models

// Success is used to indicate that the request was successful
type Success struct {
	Status int
	View   interface{}
}
