package main

import (
	"log"
	"os"

	"github.com/Fro1ko/final/pkg/db"
	"github.com/Fro1ko/final/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	server.Start()
}
