package web

//ErrorResponse how we respond to the clients when something goes wrong
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to add a web information to a request error
type Error struct {
	Err    error
	Status int
}

// NewRequestError is when a known error condition is encountered
func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
