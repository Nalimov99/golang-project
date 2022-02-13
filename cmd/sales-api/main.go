package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"garagesale/cmd/sales-api/internal/handlers"
	"garagesale/internal/platform/database"
)

func main() {
	log.Println("main : Started")
	defer log.Println("main : Completed")

	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	const timeout = 5 * time.Second

	ps := handlers.Product{
		DB: db,
	}

	api := http.Server{
		Addr:         "localhost:3020",
		Handler:      http.HandlerFunc(ps.List),
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
