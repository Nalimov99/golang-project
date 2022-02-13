package main

import (
	"flag"
	"garagesale/internal/platform/database"
	"garagesale/internal/schema"
	"log"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	var cfg struct {
		DB database.Config
	}
	err := envconfig.Process("garagesale", &cfg)
	if err != nil {
		log.Fatal(err.Error())
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
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("applying seed error: ", err)
		}

		log.Print("Seed data inserted")
		return
	default:
		log.Print("No args passed")
		return
	}
}
