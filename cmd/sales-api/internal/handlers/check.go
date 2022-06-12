package handlers

import (
	"garagesale/internal/platform/database"
	"garagesale/internal/platform/web"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Check has handlers to impement service orchestration
type Check struct {
	DB *sqlx.DB
}

// Health respond with a 200 OK if the service is healthy and ready for traffic
func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {
	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(r.Context(), c.DB); err != nil {
		health.Status = "db not ready"
		return web.Respond(r.Context(), w, health, http.StatusInternalServerError)
	}

	health.Status = "OK"
	return web.Respond(r.Context(), w, health, http.StatusOK)
}
