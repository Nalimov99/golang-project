package handlers

import (
	"garagesale/internal/middleware"
	"garagesale/internal/platform/auth"
	"garagesale/internal/platform/web"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func API(log *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {
	app := web.NewApp(log, middleware.Logger(log), middleware.Errors(log), middleware.Metric())

	c := Check{DB: db}
	app.Handle(http.MethodGet, "/v1/health", c.Health)

	u := Users{
		DB:            db,
		Log:           log,
		authenticator: authenticator,
	}
	app.Handle(http.MethodGet, "/v1/user/token", u.Token)

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
