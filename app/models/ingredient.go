package models

type Ingredient struct {
    Amount float32 `json:"amount" bson:"a" `
    Unit   string `json:"unit" bson:"u"`
    Name   string `json:"name" bson:"n"`
}

type Ingredients []Ingredient

type IngredientGroup struct {
  Title string `json:"title" bson:"r"`
  Ingredients []Ingredient `json:"ingredients" bson:"i"`
}
