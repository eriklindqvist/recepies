package app

import (
	"net/http"
	"github.com/gorilla/mux"
	c "github.com/eriklindqvist/recepies/app/controllers"
)

func NewRoutes(rc c.RecipeController) []Route {
	return Routes{
		Route{"Index", "GET",	"/", Index},

		Route{"CreateRecipe",	"POST",	  "/recipe",      CreateRecipe},
		Route{"ReadRecipe",	  "GET",    "/recipe/{id}", ReadRecipe},
		Route{"UpdateRecipe", "PUT",    "/recipe/{id}", UpdateRecipe},
		Route{"DeleteRecipe",	"DELETE", "/recipe/{id}",	DeleteRecipe},
		Route{"ListRecepies", "GET",    "/recepies",    ListRecepies},
		Route{"UploadImage",  "POST",   "/recipe/{id}/upload",  UploadImage},

		Route{"ListIngredients", "GET",    "/ingredients",    ListIngredients},
		Route{"ListRecipeNames", "GET",    "/recepies/names",    ListNames},
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

func UploadImage(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.Upload(mux.Vars(r)["id"], r)
	}, w, r)
}

func ListIngredients(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.Ingredients()
	}, w, r)
}

func ListNames(w http.ResponseWriter, r *http.Request) {
	Handle(func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
		return rc.ListNames()
	}, w, r)
}
