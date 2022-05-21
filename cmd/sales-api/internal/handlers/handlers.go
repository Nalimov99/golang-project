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

// matchPredefinedErrors knows how to respond for known failure scenarios
func matchPredefinedErrors(err error) error {
	switch err {
	case product.ErrNotFound:
		return web.NewRequestError(err, http.StatusNotFound)
	case product.ErrInvalidId:
		return web.NewRequestError(err, http.StatusBadRequest)
	default:
		return nil
	}
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
		if webErr := matchPredefinedErrors(err); webErr != nil {
			return webErr
		}

		return errors.Wrapf(err, "looking for product %v", id)
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

// UpdateProduct decodes the body of a request to update an existing Product.
func (p *Product) UpdateProduct(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var updates product.UpdateProduct
	if err := web.Decode(r, &updates); err != nil {
		return err
	}

	prod, err := product.Update(r.Context(), p.DB, id, updates, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(w, prod, http.StatusOK)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (p *Product) AddSale(w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale

	if err := web.Decode(r, &ns); err != nil {
		return err
	}

	id := chi.URLParam(r, "product_id")

	sale, err := product.AddSale(r.Context(), p.DB, ns, id, time.Now())
	if err != nil {
		if webErr := matchPredefinedErrors(err); webErr != nil {
			return webErr
		}

		return errors.Wrapf(err, "looking for product %v", id)
	}

	return web.Respond(w, sale, http.StatusCreated)
}

// ListSales gets all Sales for a particular Product
func (p *Product) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "product_id")

	sales, err := product.ListSales(r.Context(), p.DB, id)
	if err != nil {
		if webErr := matchPredefinedErrors(err); webErr != nil {
			return webErr
		}

		return errors.Wrapf(err, "looking for product %v", id)
	}

	return web.Respond(w, sales, http.StatusOK)
}
