package web

import (
	"encoding/json"
	"net/http"
)

// Decode looks for a JSON document in request body and unmarshals it into value
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	return nil
}
