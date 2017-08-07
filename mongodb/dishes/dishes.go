package dishes

import (
	"log"
	helper "mongodb/helper"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Dish is to be exported
type Dish struct {
	Type  string // type of dish, includes drinks and deserts
	Name  string // name of the dish
	Price string // preco do prato multiplicar por 100 e nao ter digits
}

// Dishadd is for export
func Dishadd(database helper.DatabaseX, dishInsert Dish) helper.Resultado {

	database.Collection = "dishes"
	dishType := dishInsert.Type
	dishName := dishInsert.Name
	dishPrice := dishInsert.Price

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Insert(&Dish{dishType, dishName, dishPrice})
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
