package dishes

import (
	"log"
	helper "mongodb/helper"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Dish is to be exported
type Dish struct {
	Type       string // type of dish, includes drinks and deserts
	Name       string // name of the dish
	Price      string // preco do prato multiplicar por 100 e nao ter digits
	GlutenFree string // Gluten free dishes
	DairyFree  string // Dairy Free dishes
	Vegetarian string // Vegeterian dishes
}

// Dishadd is for export
func Dishadd(database helper.DatabaseX, dishInsert Dish) helper.Resultado {

	database.Collection = "dishes"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Insert(dishInsert)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Something Happened"
	res.IsSuccessful = "Y"

	return res
}

func find(database helper.DatabaseX, dishFind Dish) {

	database.Collection = "dishes"

	dishName := dishFind.Name

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	result := Dish{}
	err = c.Find(bson.M{"Name": dishName}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

}

// GetAll works
func GetAll(database helper.DatabaseX) []Dish {

	database.Collection = "dishes"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	var results []Dish

	err = c.Find(nil).All(&results)
	if err != nil {
		// TODO: Do something about the error
	} else {
		return results
	}

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
