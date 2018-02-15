
package controllers

import (
	"net/http"
	"io"
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"mime/multipart"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"image/gif"
	"github.com/eriklindqvist/recepies_auth/log"
	l "github.com/eriklindqvist/recepies_api/app/lib"
	m "github.com/eriklindqvist/recepies_api/app/models"
)

type RecipeController struct {
		db *mgo.Database
		c *mgo.Collection
		elastic string
}

func (rc RecipeController) Elastic(method string, id string, data io.Reader) (resp *http.Response, err error) {
	if (rc.elastic != "") {
		c := http.Client{}
		req, err := http.NewRequest(method, rc.elastic+"recepies/recipe/"+id, data)

		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		return c.Do(req)
	}
	return
}

func NewRecipeController(s *mgo.Session) *RecipeController {
		db := s.DB(l.Getenv("DATABASE", "recepies"))
		e := os.Getenv("ELASTIC_HOST")
		elastic := ""

		if (e != "") {
			p := os.Getenv("ELASTIC_PORT")
			elastic = fmt.Sprintf("http://%s:%s/", e, p)
		}

    return &RecipeController{db, db.C("recepies"), elastic}
}

func (rc RecipeController) Create(json io.Reader) ([]byte, error) {
		r := new(m.Recipe)

		if err := r.FromJson(json); err != nil {
				log.Err(fmt.Sprintf("Error parsing JSON: %s", err.Error()))
				return nil, l.NewError(http.StatusBadRequest, err.Error())
		}

		if err := r.Insert(rc.c); err != nil {
				log.Err(fmt.Sprintf("Error inserting document: %s", err.Error()))
				return nil, err
		}

		if _, err := rc.Elastic("PUT", r.Id.Hex(), json); err != nil {
			log.Err(fmt.Sprintf("Error updating elstic: %s", err.Error()))
		}

		return r.ToJson()
}

func (rc RecipeController) Read(id string) ([]byte, error) {
		if !bson.IsObjectIdHex(id) {
				log.Err(fmt.Sprintf("Invalid BSON ID: %s", id))
        return nil, l.NewError(http.StatusNotFound, "Recipe not found")
    }

		recipe := new(m.Recipe)

		if err := recipe.Find(bson.ObjectIdHex(id), rc.c); err != nil {
				log.Err(fmt.Sprintf("Error reading recipe %s: %s", id, err.Error()))
        return nil, l.NewError(http.StatusNotFound, "Recipe not found")
    }

		return recipe.ToJson()
}

func (rc RecipeController) Update(id string, json io.Reader) ([]byte, error) {
		if !bson.IsObjectIdHex(id) {
				log.Err(fmt.Sprintf("Invalid BSON ID: %s", id))
				return nil, l.NewError(http.StatusNotFound, "Recipe not found")
		}

		r := new(m.Recipe)

		if err := r.FromJson(json); err != nil {
				log.Err(fmt.Sprintf("Error parsing JSON: %s", err.Error()))
				return nil, l.NewError(http.StatusBadRequest, err.Error())
		}

		if (bson.ObjectIdHex(id) != r.Id) {
				return nil, l.NewError(http.StatusBadRequest, "ID parameter must equal ID in body")
		}

		if err := r.Update(rc.c); err != nil {
				if (err.Error() == "not found") {
					err = l.NewError(http.StatusNotFound, "Recipe not found")
				}

				log.Err(fmt.Sprintf("Error updating recipe: %s", err.Error()))
				return nil, err
		}

		if _, err := rc.Elastic("PUT", r.Id.Hex(), json); err != nil {
			log.Err(fmt.Sprintf("Error updating elstic: %s", err.Error()))
		}

		return r.ToJson()
}


func (rc RecipeController) Delete(id string) ([]byte, error) {
		if !bson.IsObjectIdHex(id) {
				log.Err(fmt.Sprintf("Invalid BSON ID: %s", id))
        return nil, l.NewError(http.StatusNotFound, "Recipe not found")
    }

		r := m.Recipe{Id: bson.ObjectIdHex(id)}

		err := r.Delete(rc.c)

		if err != nil {
			if err.Error() == "not found" {
				err = l.NewError(http.StatusNotFound, "Recipe not found")
			}

			log.Err(fmt.Sprintf("Error deleting recipe: %s", err.Error()))
		}

		if _, err := rc.Elastic("DELETE", r.Id.Hex(), nil); err != nil {
			log.Err(fmt.Sprintf("Error deleting from elstic: %s", err.Error()))
		}

		return nil, err
}

func (rc RecipeController) List() ([]byte, error) {
		r := m.Recepies{}

		if err := r.List(rc.c); err != nil {
			log.Err(fmt.Sprintf("Error listing recepies: %s", err.Error()))
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
		ok bool
	)

	if !bson.IsObjectIdHex(id) {
			log.Err(fmt.Sprintf("Invalid BSON ID: %s", id))
			return nil, l.NewError(http.StatusNotFound, "Recipe not found")
	}

	recipe := new(m.Recipe)

	if err := recipe.Find(bson.ObjectIdHex(id), rc.c); err != nil {
			log.Err(fmt.Sprintf("Error reading recipe %s: %s", id, err.Error()))
			return nil, l.NewError(http.StatusNotFound, "Recipe not found")
	}

	if reader, err = r.MultipartReader(); err != nil {
		log.Err("Could not get multipart reader")
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

			name := part.FileName()

			if _, ok = ctypes[ct]; !ok {
				return nil, l.NewError(http.StatusUnsupportedMediaType, "Illegal content type")
			}

			filename := l.Getenv("FILEBASE", "/files") + "/" + name

			if dst, err = os.Create(filename); err != nil {
					log.Err("Error creating file:")
          return nil, l.NewError(http.StatusInternalServerError, err.Error())
      }

			defer dst.Close()

			if _, err = io.Copy(dst, part); err != nil {
				log.Err("Error saving file:")
				return nil, l.NewError(http.StatusInternalServerError, err.Error())
			}

			file, err := os.Open(filename)

			if err != nil {
				log.Err("Error opening file:")
				return nil, l.NewError(http.StatusInternalServerError, err.Error())
			}

			var img image.Image

			img, _, err = image.Decode(file)

			if err != nil {
				log.Err("Error decoding file:")
				return nil, l.NewError(http.StatusInternalServerError, err.Error())
			}

			file.Close()

			m := resize.Thumbnail(320, 240, img, resize.Lanczos3)

			thumbsdir := l.Getenv("FILEBASE", "/files") + "/thumbs"

			if _, err := os.Stat(thumbsdir); err != nil {
    		if os.IsNotExist(err) {
        	os.Mkdir(thumbsdir, 0777) //TODO: Move this check to system startup instead
			  }
			}

			out, err := os.Create(thumbsdir + "/" + name)

			if err != nil {
				log.Err("Error creating thumbnail:")
				return nil, l.NewError(http.StatusInternalServerError, err.Error())
			}

			defer out.Close()

			if (ct == "image/jpeg") {
				jpeg.Encode(out, m, nil)
			} else if (ct == "image/gif") {
				gif.Encode(out, m, nil)
			} else if (ct == "image/png") {
				png.Encode(out, m)
			}

			recipe.Image = name
		}

	err = recipe.Update(rc.c)

	return nil,err
}

func (rc RecipeController) Ingredients() ([]byte, error) {
		var names []string

		if err := rc.c.Find(nil).Distinct("j.i.n", &names); err != nil {
			log.Err("Error creating file:")
			err = l.NewError(http.StatusInternalServerError, err.Error())
		}

		return json.Marshal(names)
}

func (rc RecipeController) Units() ([]byte, error) {
		var units []string

		if err := rc.c.Find(nil).Distinct("j.i.u", &units); err != nil {
			log.Err("Error listing units:")
			err = l.NewError(http.StatusInternalServerError, err.Error())
		}

		return json.Marshal(units)
}

type Names []struct{ Id  bson.ObjectId `json:"id" bson:"_id"`; Title string `json:"title" bson:"t"` }

func (rc RecipeController) ListNames() ([]byte, error) {
	r := Names{}

	if err := rc.c.Find(nil).Select(bson.M{"t": 1}).Limit(100).All(&r); err != nil {
		log.Err("Error listing names:")
		return nil, err
	}

	return json.Marshal(r)
}
