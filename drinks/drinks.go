package drinks

import (
	"log"

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
func Drinkadd(drinkinsert Drinks) string {

	drinktype := drinkinsert.Type
	drinkname := drinkinsert.Name
	drinkprice := drinkinsert.Price

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("restaurante").C("dish")

	result := Drinks{}
	err = c.Find(bson.M{"name": drinkname}).One(&result)

	if err != nil {
		if err.Error() == "not found" {
			err = c.Insert(&Drinks{drinktype, drinkname, drinkprice})
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
