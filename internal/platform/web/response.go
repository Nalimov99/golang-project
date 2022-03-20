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
