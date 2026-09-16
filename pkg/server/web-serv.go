package server

import (
	"log"
	"net/http"
	"os"

	"github.com/Fro1ko/final/pkg/api"
)

func Start() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Printf("Starting web server on port %s", port)

	return http.ListenAndServe(":"+port, nil)
}
