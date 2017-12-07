package models

type Ingredient struct {
    Amount uint16 `json:"amount" bson:"a" `
    Unit   string `json:"unit" bson:"u"`
    Name   string `json:"name" bson:"n"`
}

type Ingredients []Ingredient
