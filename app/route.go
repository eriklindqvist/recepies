package app

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	l "github.com/eriklindqvist/recepies/app/lib"
	c "github.com/eriklindqvist/recepies/app/controllers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route
type Endpoint func(w http.ResponseWriter, r *http.Request) ([]byte, error)

func getSession() *mgo.Session {
		host := "mongodb://" + l.Getenv("MONGODB_HOST", "localhost")
    s, err := mgo.Dial(host)
		log.Printf("host: %s", host)
    // Check if connection error, is mongo running?
    if err != nil {

        panic(err)
    }
    return s
}

var rc = c.NewRecipeController(getSession())
var routes = NewRoutes(*rc)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Handle(endpoint Endpoint, w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	body, err := endpoint(w, r)

	if err != nil {
		switch e := err.(type) {
		case l.Error:
			http.Error(w, e.Error(), e.Status())
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	writeBody(w, body)
}

func setContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func writeNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func writeBody(w http.ResponseWriter, body []byte) {
	fmt.Fprintf(w, "%s", body)
}
