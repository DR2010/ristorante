// Package btcmarkets is a dish for packages
// -------------------------------------
// .../restauranteapi/dishes/btcmarkets.go
// -------------------------------------
package btcmarkets

import (
	"fmt"
	"log"

	helper "restauranteapi/helper"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"
)

// PreOrder is to be exported
type PreOrder struct {
	Currency string // Currency
	Max      string // balance
	Min      string // Cotacao
	Email    string // date time
	Buy      string // date time
	Sell     string // date time
	Date     string
	DateTime string
}

// DCPreOrder handling preorders
type DCPreOrder struct {
	Preorders []PreOrder
}

// PreOrderAdd is for export
func PreOrderAdd(redisclient *redis.Client, cryptoInsert DCPreOrder) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "btcpreorder"
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

	err = collection.Insert(cryptoInsert)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Crypto Cotacao added"
	res.IsSuccessful = "Y"

	return res
}

// PreorderGetAll works
func PreorderGetAll(redisclient *redis.Client) []PreOrder {

	database := new(helper.DatabaseX)

	database.Collection = "btcpreorder"

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

	var results []PreOrder

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
