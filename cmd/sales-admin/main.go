package main

import (
	"flag"
	"garagesale/internal/platform/database"
	"garagesale/internal/schema"
	"log"
)

func main() {
	db, err := database.Open()
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
