package handlers

import (
	"encoding/json"
	"garagesale/internal/product"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Product struct {
	DB *sqlx.DB
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error querying db ", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error marshaling ", err)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		log.Println(err)
	}
}
