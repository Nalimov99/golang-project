package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("main : Started")
	defer log.Println("main : Completed")

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
