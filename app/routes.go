package app

import (
	"net/http"
	"github.com/gorilla/mux"
	c "./controllers"
)

func NewRoutes(rc c.RecipeController) []Route {
	return Routes{
		Route{"Index", "GET",	"/", Index},

		Route{"CreateRecipe",	"POST",	  "/recipe",      CreateRecipe},
		Route{"ReadRecipe",	  "GET",    "/recipe/{id}", ReadRecipe},
		Route{"UpdateRecipe", "PUT",    "/recipe/{id}", UpdateRecipe},
		Route{"DeleteRecipe",	"DELETE", "/recipe/{id}",	DeleteRecipe},
		Route{"ListRecepies", "GET",    "/recepies",    ListRecepies},
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return []byte{}, nil
	}, w, r)
}

func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.Create(r.Body)
	}, w, r)
}

func ReadRecipe(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.Read(mux.Vars(r)["id"])
	}, w, r)
}

func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.Update(mux.Vars(r)["id"], r.Body)
	}, w, r)
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.Delete(mux.Vars(r)["id"])
	}, w, r)
}

func ListRecepies(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.List()
	}, w, r)
}
