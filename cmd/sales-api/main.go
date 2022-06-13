package main

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "expvar" // register the /debug/vars handlers
	"garagesale/cmd/sales-api/internal/handlers"
	"garagesale/internal/platform/auth"
	"garagesale/internal/platform/database"
	_ "net/http/pprof" // Register the /debug/pprof handlers

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// =======================================================
	// Setup dependencies

	log.Println("main : Started")
	defer log.Println("main : Completed")

	log := log.New(os.Stdout, "SALES: ", log.LstdFlags|log.Lshortfile)

	var cfg struct {
		DB     database.Config
		Server struct {
			Addr                  string        `default:"localhost:3020"`
			Debug                 string        `default:"localhost:6060"`
			ReadTimeout           time.Duration `default:"5s" split_words:"true"`
			WriteTimeout          time.Duration `default:"5s" split_words:"true"`
			GracefullShutdownTime time.Duration `default:"5s" split_words:"true"`
		}
		Auth struct {
			KeyID              string `default:"1"`
			PrivateKeyFromFile string `default:"private.pem"`
			Algorithm          string `default:"RS256"`
		}
	}
	err := envconfig.Process("garagesale", &cfg)
	if err != nil {
		return errors.Wrap(err, "generating config usage")
	}

	const dbConfigFormat = "\n\nDatabse config\nUser: %v\nPassword: %v\nHost: %v\nPath: %v\nSslMode: %v\n\n"
	log.Printf(dbConfigFormat, cfg.DB.Host, cfg.DB.Password, cfg.DB.Host, cfg.DB.Path, cfg.DB.SslMode)

	// =======================================================
	// Initialize authentication support

	authenticator, err := createAuth(
		cfg.Auth.PrivateKeyFromFile,
		cfg.Auth.KeyID,
		cfg.Auth.Algorithm,
	)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}

	// =======================================================
	// Open DB

	db, err := database.Open(cfg.DB)
	if err != nil {
		return errors.Wrap(err, "Opening db")
	}
	defer db.Close()

	// =======================================================
	// Start debug service

	go func() {
		log.Printf("main : Debug listen on %s", cfg.Server.Debug)
		http.ListenAndServe(cfg.Server.Debug, nil)
	}()

	// =======================================================
	// Start API service

	api := http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      handlers.API(log, db, authenticator),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.ReadTimeout,
	}
	const serverConfigFormat = "\n\nServer config:\nAddress: %v\nReadTimeout: %v\nWriteTimeout: %v\nGracefullShutdown: %v\n\n"
	log.Printf(serverConfigFormat, cfg.Server.Addr, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout, cfg.Server.GracefullShutdownTime)

	// =======================================================
	// Listen to shutdown server

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

func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {
	keyContent, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContent)
	if err != nil {
		return nil, errors.Wrap(err, "parsing private key")
	}

	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, keyID, algorithm, public)
}
