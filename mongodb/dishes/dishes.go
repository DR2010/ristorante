package dishes

import (
	"fmt"
	"log"
	helper "restauranteapi/helper"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Dish is to be exported
type Dish struct {
	SystemID   bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Name       string        // name of the dish - this is the KEY, must be unique
	Type       string        // type of dish, includes drinks and deserts
	Price      string        // preco do prato multiplicar por 100 e nao ter digits
	GlutenFree string        // Gluten free dishes
	DairyFree  string        // Dairy Free dishes
	Vegetarian string        // Vegeterian dishes
}

// Dishadd is for export
func Dishadd(redisclient *redis.Client, dishInsert Dish) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "dishes"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

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
	res.ErrorDescription = "Dish added"
	res.IsSuccessful = "Y"

	return res
}

// Find is to find stuff
func Find(redisclient *redis.Client, dishFind string) (Dish, string) {

	database := new(helper.DatabaseX)
	database.Collection = "dishes"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	dishName := dishFind
	dishnull := Dish{}

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	result := []Dish{}
	err1 := c.Find(bson.M{"name": dishName}).All(&result)
	if err1 != nil {
		log.Fatal(err1)
	}

	var numrecsel = len(result)

	if numrecsel <= 0 {
		return dishnull, "404 Not found"
	}

	return result[0], "200 OK"
}

// GetAll works
func GetAll(redisclient *redis.Client) []Dish {

	database := new(helper.DatabaseX)

	database.Collection = "dishes"

	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	fmt.Println("database.Location")
	fmt.Println(database.Location)

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

// Dishupdate is
func Dishupdate(redisclient *redis.Client, dishUpdate Dish) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "dishes"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Update(bson.M{"name": dishUpdate.Name}, dishUpdate)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Something Happened"
	res.IsSuccessful = "Y"

	return res
}

// Dishdelete is
func Dishdelete(redisclient *redis.Client, dishDelete Dish) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "dishes"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()
	database.Collection = "dishes"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Remove(bson.M{"name": dishDelete.Name})

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Dish deleted successfully"
	res.IsSuccessful = "Y"

	return res
}
