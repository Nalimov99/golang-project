package handlers

import (
	"garagesale/internal/platform/web"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func API(log *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(log)

	p := Product{
		DB:  db,
		Log: log,
	}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	return app
}
