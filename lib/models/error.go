package models

// RequestError is an error wrapper type that wraps request's error as
// well as it's http status code
type RequestError struct {
	Status int
	Err    error
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}
