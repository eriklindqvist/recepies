package controllers

import (
	"net/http"
	"io"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	l "../lib"
	m "../models"
)

type RecipeController struct {
		db *mgo.Database
		c *mgo.Collection
}

func NewRecipeController(s *mgo.Session) *RecipeController {
		db := s.DB(l.Getenv("DATABASE", "recepies"))

    return &RecipeController{db, db.C("recepies")}
}

func (rc RecipeController) Create(json io.Reader) ([]byte, error) {
		r := new(m.Recipe)

		if err := r.FromJson(json); err != nil {
				return nil, l.NewError(http.StatusBadRequest, err.Error())
		}

		if err := r.Insert(rc.c); err != nil {
				return nil, err
		}

		return r.ToJson()
}

func (rc RecipeController) Read(id string) ([]byte, error) {
		if !bson.IsObjectIdHex(id) {
        return nil, l.NewError(http.StatusNotFound, "Recipe not found")
    }

		recipe := new(m.Recipe)

		if err := recipe.Find(bson.ObjectIdHex(id), rc.c); err != nil {
        return nil, l.NewError(http.StatusNotFound, "Recipe not found")
    }

		return recipe.ToJson()
}

func (rc RecipeController) Update(id string, json io.Reader) ([]byte, error) {
		if !bson.IsObjectIdHex(id) {
				return nil, l.NewError(http.StatusNotFound, "Recipe not found")
		}

		r := new(m.Recipe)

		if err := r.FromJson(json); err != nil {
				return nil, l.NewError(http.StatusBadRequest, err.Error())
		}

		if (bson.ObjectIdHex(id) != r.Id) {
				return nil, l.NewError(http.StatusBadRequest, "ID parameter must equal ID in body")
		}

		if err := r.Update(rc.c); err != nil {
				if (err.Error() == "not found") {
					err = l.NewError(http.StatusNotFound, "Recipe not found")
				}

				return nil, err
		}

		return r.ToJson()
}


func (rc RecipeController) Delete(id string) ([]byte, error) {
		if !bson.IsObjectIdHex(id) {
        return nil, l.NewError(http.StatusNotFound, "Recipe not found")
    }

		r := m.Recipe{Id: bson.ObjectIdHex(id)}

		err := r.Delete(rc.c)

		if err != nil && err.Error() == "not found" {
				err = l.NewError(http.StatusNotFound, "Recipe not found")
		}

		return nil, err
}

func (rc RecipeController) List() ([]byte, error) {
		r := m.Recepies{}

		if err := r.List(rc.c); err != nil {
			return nil, err
		}

		return json.Marshal(r)
}

func (rc RecipeController) Ingredients() ([]byte, error) {
		var names []string

		if err := rc.c.Find(nil).Distinct("i.n", &names); err != nil {
			err = l.NewError(http.StatusInternalServerError, err.Error())
		}

		return json.Marshal(names)
}
