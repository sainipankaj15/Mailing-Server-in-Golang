package main

import (
	"database/sql"
	"log"
	jsonapi "mailingServer/jsonAPI"
	"mailingServer/mailDatabase"
	"sync"

	"github.com/alexflint/go-arg"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
}

func main() {

	arg.MustParse(&args)

	if args.DbPath == "" {
		args.DbPath = "list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":8080"
	}

	log.Printf("Using Database : '%v'\n", args.DbPath)

	db, err := sql.Open("sqlite3", args.DbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	mailDatabase.TryCreate(db)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		log.Println("Starting JSON API Server .......")
		jsonapi.Serve(db, args.BindJson)
	}()

	wg.Wait()

}
