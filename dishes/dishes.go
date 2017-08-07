package dishes

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Dish struct {
	Type  string // type of dish, includes drinks and deserts
	Name  string // name of the dish
	Price int32  // preco do prato multiplicar por 100 e nao ter digits
}

func add() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("restaurante").C("people")

	err = c.Insert(&Dish{"Entree", "Pao de Queijo", 1000}, &Dish{"Main", "Feijoada", 5000})
	if err != nil {
		log.Fatal(err)
	}

	result := Dish{}
	err = c.Find(bson.M{"Name": "Feijoada"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Price:", result.Price)
}
