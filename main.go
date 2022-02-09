package main

import (
	"context"
	"encoding/json"
	"flag"
	"garagesale/schema"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("main : Started")
	defer log.Println("main : Completed")

	db, err := openDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("applying migrations error: ", err)
		}

		log.Print("Migrations complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("applying seed error: ", err)
		}

		log.Print("Seed data inserted")
		return
	}

	const timeout = 5 * time.Second

	api := http.Server{
		Addr:         "localhost:3020",
		Handler:      http.HandlerFunc(ListProducts),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error : listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Gracefull shutdown not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : Could not stop server gracefully : %v", err)
		}
	}
}

func openDb() (*sqlx.DB, error) {
	q := url.Values{}

	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")
	q.Set("port", "5434")
	log.Print(q.Encode())
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "1234"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	log.Print(u.String())
	log.Print(u.String())
	return sqlx.Open("postgres", u.String())
}

type Products struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Cost     int    `json:"cost"`
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	list := []Products{
		{Name: "Book", Quantity: 2, Cost: 10},
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
