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
		args.DbPath = "mailingList.db"
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

	// Check if the connection is alive by executing a test query
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	log.Println("Connected to the database.")

	mailDatabase.TryCreate(db)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Println("Starting JSON API Server .......")
		jsonapi.Serve(db, args.BindJson)
	}()

	wg.Wait()

}
