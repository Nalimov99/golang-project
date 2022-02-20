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

	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("main : Started")
	defer log.Println("main : Completed")

	log := log.New(os.Stdout, "SALES: ", log.LstdFlags|log.Lshortfile)

	var cfg struct {
		DB     database.Config
		Server struct {
			Addr                  string        `default:"localhost:3020"`
			ReadTimeout           time.Duration `default:"5s" split_words:"true"`
			WriteTimeout          time.Duration `default:"5s" split_words:"true"`
			GracefullShutdownTime time.Duration `default:"5s" split_words:"true"`
		}
	}
	err := envconfig.Process("garagesale", &cfg)
	if err != nil {
		return errors.Wrap(err, "generating config usage")
	}

	const dbConfigFormat = "\n\nDatabse config\nUser: %v\nPassword: %v\nHost: %v\nPath: %v\nSslMode: %v\n\n"
	log.Printf(dbConfigFormat, cfg.DB.Host, cfg.DB.Password, cfg.DB.Host, cfg.DB.Path, cfg.DB.SslMode)

	db, err := database.Open(cfg.DB)
	if err != nil {
		return errors.Wrap(err, "Opening db")
	}
	defer db.Close()

	api := http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      handlers.API(log, db),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.ReadTimeout,
	}
	const serverConfigFormat = "\n\nServer config:\nAddress: %v\nReadTimeout: %v\nWriteTimeout: %v\nGracefullShutdown: %v\n\n"
	log.Printf(serverConfigFormat, cfg.Server.Addr, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout, cfg.Server.GracefullShutdownTime)

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "listening and serving")

	case <-shutdown:
		log.Println("main : Start shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefullShutdownTime)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Gracefull shutdown not complete in %v : %v", cfg.Server.GracefullShutdownTime, err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "gracefull shutdown")
		}
	}

	return nil
}
