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

// ServerError is an error wrapper type that wraps server's error as
// well as it's http status code
type ServerError struct {
	Status int
	Reason error
}

func (se *ServerError) Error() string {
	return se.Reason.Error()
}
