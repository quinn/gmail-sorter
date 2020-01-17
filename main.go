package main

import (
	"github.com/quinn/gmail-sorter/pkg/core"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
)

func main() {
	var err error
	api, err := gmailapi.Start()

	if err != nil {
		panic(err)
	}

	db := db.NewDB()
	defer db.Close()

	spec, err := core.NewSpec(api, db)

	if err != nil {
		panic(err)
	}

	err = spec.Apply()

	if err != nil {
		panic(err)
	}
}
