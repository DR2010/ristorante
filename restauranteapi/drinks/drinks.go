package drinks

import (
	"log"
	"restauranteapi/helper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Drinks is a struct
type Drinks struct {
	Type  string // type of dish, includes drinks and deserts
	Name  string // name of the dish
	Price string // preco do prato multiplicar por 100 e nao ter digits
}

// Drinkadd does the job
func Drinkadd(database helper.DatabaseX, drinkinsert Drinks) string {

	database.Collection = "drinks"
	drinktype := drinkinsert.Type
	drinkname := drinkinsert.Name
	drinkprice := drinkinsert.Price

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	result := Drinks{}
	err = collection.Find(bson.M{"name": drinkname}).One(&result)

	if err != nil {
		if err.Error() == "not found" {
			err = collection.Insert(&Drinks{drinktype, drinkname, drinkprice})
			if err != nil {
				log.Fatal(err)
				return err.Error()
			}

			return "Dish created"

		}
		return "something went wrong"
	}

	return "Drink already exists"
}
