package handlers

import (
	"context"
	"garagesale/internal/platform/auth"
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
	case product.ErrForbidden:
		return web.NewRequestError(err, http.StatusForbidden)
	default:
		return nil
	}
}

// List gives all known products
func (p *Product) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(ctx, p.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Retrieve gives a single product
func (p *Product) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(ctx, p.DB, id)
	if err != nil {
		if webErr := matchPredefinedErrors(err); webErr != nil {
			return webErr
		}

		return errors.Wrapf(err, "looking for product %v", id)
	}

	return web.Respond(ctx, w, prod, http.StatusOK)
}

//Create decode a JSON from a POST request and create a new product
func (p *Product) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("auth claims not in context")
	}

	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(ctx, p.DB, claims, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, prod, http.StatusCreated)
}

// UpdateProduct decodes the body of a request to update an existing Product.
func (p *Product) UpdateProduct(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("auth claims not in context")
	}

	id := chi.URLParam(r, "id")

	var updates product.UpdateProduct
	if err := web.Decode(r, &updates); err != nil {
		return err
	}

	prod, err := product.Update(ctx, p.DB, claims, id, updates, time.Now())
	if err != nil {
		if webErr := matchPredefinedErrors(err); webErr != nil {
			return webErr
		}

		return errors.Wrapf(err, "updating product %v", id)
	}

	return web.Respond(ctx, w, prod, http.StatusOK)
}

// DeleteProduct removes a single Product indentified by an ID in the request URL
func (p *Product) DeleteProduct(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := product.Delete(ctx, p.DB, id); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (p *Product) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale

	if err := web.Decode(r, &ns); err != nil {
		return err
	}

	id := chi.URLParam(r, "product_id")

	sale, err := product.AddSale(ctx, p.DB, ns, id, time.Now())
	if err != nil {
		if webErr := matchPredefinedErrors(err); webErr != nil {
			return webErr
		}

		return errors.Wrapf(err, "looking for product %v", id)
	}

	return web.Respond(ctx, w, sale, http.StatusCreated)
}

// ListSales gets all Sales for a particular Product
func (p *Product) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "product_id")

	sales, err := product.ListSales(ctx, p.DB, id)
	if err != nil {
		if webErr := matchPredefinedErrors(err); webErr != nil {
			return webErr
		}

		return errors.Wrapf(err, "looking for product %v", id)
	}

	return web.Respond(ctx, w, sales, http.StatusOK)
}
