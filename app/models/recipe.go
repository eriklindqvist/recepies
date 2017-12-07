package models

import (
  "io"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Recipe struct {
    Id     bson.ObjectId `json:"id" bson:"_id"`
    Title   string `json:"title" bson:"t"`
    Description string `json:"description" bson:"d"`
    Ingredients []Ingredient `json:"ingredients" bson:"i"`
}

type Recepies []Recipe

func (r *Recipe) Find(id bson.ObjectId, c *mgo.Collection) error {
  return c.FindId(id).One(&r)
}

func (r *Recipe) FromJson(reader io.Reader) error {
  return json.NewDecoder(reader).Decode(r)
}

func (r Recipe) ToJson() ([]byte, error) {
  return json.Marshal(r)
}

func (r *Recipe) Insert(c *mgo.Collection) error {
  r.Id = bson.NewObjectId()
  return c.Insert(r)
}

func (r *Recipe) Update(c *mgo.Collection) error {
  return c.UpdateId(r.Id, r)
}

func (r *Recipe) Delete(c *mgo.Collection) error {
  return c.RemoveId(r.Id)
}

func (r *Recepies) List(c *mgo.Collection) error {
  return c.Find(nil).Limit(100).All(r)
}
