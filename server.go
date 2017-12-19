package main

import (
	"github.com/eriklindqvist/recepies/app"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	router := app.NewRouter()

	log.Fatal(http.ListenAndServe(":3000", router))
}
