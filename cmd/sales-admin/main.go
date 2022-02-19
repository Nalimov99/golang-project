package main

import (
	"flag"
	"garagesale/internal/platform/database"
	"garagesale/internal/schema"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cfg struct {
		DB database.Config
	}
	err := envconfig.Process("garagesale", &cfg)
	if err != nil {
		return errors.Wrap(err, "generating config usage")
	}

	db, err := database.Open(cfg.DB)
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
		return nil
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("applying seed error: ", err)
		}

		log.Print("Seed data inserted")
		return nil
	default:
		log.Print("No args passed")
		return nil
	}
}
