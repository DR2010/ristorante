// Package models is a dish for packages
// -------------------------------------
// .../restauranteapi/models/dishes.go
// -------------------------------------
package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Dish is to be exported
type Dish struct {
	SystemID         bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Name             string        // name of the dish - this is the KEY, must be unique
	Type             string        // type of dish, includes drinks and deserts
	Price            string        // preco do prato multiplicar por 100 e nao ter digits
	GlutenFree       string        // Gluten free dishes
	DairyFree        string        // Dairy Free dishes
	Vegetarian       string        // Vegeterian dishes
	InitialAvailable string        // Number of items initially available
	CurrentAvailable string        // Currently available
	ImageName        string        // Image Name
	Description      string        // Description of the plate
	Descricao        string        // Descricao do prato

}
