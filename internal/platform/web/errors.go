package web

//ErrorResponse how we respond to the clients when something goes wrong
type ErrorResponse struct {
	Error string `json:"error"`
}
