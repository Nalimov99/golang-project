package handlers

import (
	"encoding/json"
	"garagesale/internal/product"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db ", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshaling ", err)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		p.Log.Println(err)
	}
}

func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db ", err)
		return
	}

	data, err := json.Marshal(prod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshaling ", err)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Log.Println(err)
	}
}
