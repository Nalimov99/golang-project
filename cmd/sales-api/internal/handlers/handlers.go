package handlers

import (
	"garagesale/internal/platform/web"
	"garagesale/internal/product"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gives all known products
func (p *Product) List(w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(r.Context(), p.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrieve gives a single product
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(r.Context(), p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for product %q", id)
		}
	}

	return web.Respond(w, prod, http.StatusOK)
}

//Create decode a JSON from a POST request and create a new product
func (p *Product) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(r.Context(), p.DB, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(w, prod, http.StatusCreated)
}
