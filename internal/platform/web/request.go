package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Decode looks for a JSON document in request body and unmarshals it into value
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return errors.Wrap(err, "decoding request body")
	}

	return nil
}
