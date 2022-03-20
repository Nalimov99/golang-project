package handlers

import (
	"database/sql"
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
func (p *Product) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db ", err)
		return
	}

	if err := web.Respond(w, list, http.StatusOK); err != nil {
		p.Log.Println(err, "error responding")
	}
}

// Retrieve gives a single product
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			p.Log.Printf("%v. id: %v", err, id)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db ", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusOK); err != nil {
		p.Log.Println(err, "error responding")
	}
}

//Create decode a JSON from a POST request and create a new product
func (p *Product) Create(w http.ResponseWriter, r *http.Request) {
	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db ", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusCreated); err != nil {
		p.Log.Println(err, "error responding")
	}
}
