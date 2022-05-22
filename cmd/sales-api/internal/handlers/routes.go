package handlers

import (
	"garagesale/internal/platform/web"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func API(log *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(log)

	c := Check{DB: db}
	app.Handle(http.MethodGet, "/v1/health", c.Health)

	p := Product{
		DB:  db,
		Log: log,
	}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
	app.Handle(http.MethodPatch, "/v1/products/{id}", p.UpdateProduct)
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.DeleteProduct)

	app.Handle(http.MethodPost, "/v1/products/{product_id}/sales", p.AddSale)
	app.Handle(http.MethodGet, "/v1/products/{product_id}/sales", p.ListSales)

	return app
}
