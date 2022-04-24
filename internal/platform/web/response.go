package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond marshals value to a JSON and send it to the client
func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {
	data, err := json.Marshal(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.Wrap(err, "error marshalling")
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "writing to client")
	}

	return nil
}

// Respond error knows how to handle errors going out to the client
func RespondError(w http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok {
		resp := ErrorResponse{
			Error: webErr.Err.Error(),
		}

		return Respond(w, resp, webErr.Status)
	}

	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(w, resp, http.StatusInternalServerError)
}
