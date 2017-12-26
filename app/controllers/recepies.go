package controllers

import (
	"net/http"
	"io"
	"os"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"mime/multipart"
	l "github.com/eriklindqvist/recepies/app/lib"
	m "github.com/eriklindqvist/recepies/app/models"
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

func (rc RecipeController) Upload(id string, r *http.Request) ([]byte, error) {
	var (
		reader *multipart.Reader
	 	part *multipart.Part
	 	dst *os.File
	 	err error
		ending string
		ok bool
	)

	if !bson.IsObjectIdHex(id) {
			return nil, l.NewError(http.StatusNotFound, "Recipe not found")
	}

	recipe := new(m.Recipe)

	if err := recipe.Find(bson.ObjectIdHex(id), rc.c); err != nil {
			return nil, l.NewError(http.StatusNotFound, "Recipe not found")
	}

	if reader, err = r.MultipartReader(); err != nil {
		return nil, l.NewError(http.StatusInternalServerError, err.Error())
	}

	ctypes := map[string]string {
		"image/jpeg":".jpg",
		"image/gif":".gif",
		"image/png":".png",
	}

	for {
			if part, err = reader.NextPart(); err == io.EOF {
				break
			}

			ct := part.Header.Get("Content-Type")

			if ending, ok = ctypes[ct]; !ok {
				return nil, l.NewError(http.StatusUnsupportedMediaType, "Illegal content type")
			}

			filename := l.Getenv("FILEBASE", "/home/erik") + "/" + id + ending

			if dst, err = os.Create(filename); err != nil {
          return nil, l.NewError(http.StatusInternalServerError, err.Error())
      }

			defer dst.Close()

			if _, err = io.Copy(dst, part); err != nil {
				return nil, l.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil,nil
}

func (rc RecipeController) Ingredients() ([]byte, error) {
		var names []string

		if err := rc.c.Find(nil).Distinct("i.n", &names); err != nil {
			err = l.NewError(http.StatusInternalServerError, err.Error())
		}

		return json.Marshal(names)
}

type Names []struct{ Id  bson.ObjectId `json:"id" bson:"_id"`; Title string `json:"title" bson:"t"` }

func (rc RecipeController) ListNames() ([]byte, error) {
	r := Names{}
	
	if err := rc.c.Find(nil).Select(bson.M{"t": 1}).Limit(100).All(&r); err != nil {
		return nil, err
	}

	return json.Marshal(r)
}
