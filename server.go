package main

import (
	"github.com/eriklindqvist/recepies_api/app"
	"net/http"
	"github.com/eriklindqvist/recepies_auth/log"
)

func main() {
	log.Info("Server started")

	router := app.NewRouter()

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Panic(err.Error())
	}
}
