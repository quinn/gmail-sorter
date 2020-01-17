package main

import (
	"github.com/quinn/gmail-sorter/pkg/core"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting main")

	var err error
	api, err := gmailapi.Start()

	if err != nil {
		panic(err)
	}

	db := db.NewDB()
	log.Info("connected to database")

	defer func() {
		log.Info("closing database")
		err = db.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	spec, err := core.NewSpec(api, db)

	if err != nil {
		log.Fatal(err)
	}

	err = spec.Apply()

	if err != nil {
		log.Fatal(err)
	}
}
