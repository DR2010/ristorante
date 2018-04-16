// Package btcmarkets is a dish for packages
// -------------------------------------
// .../restauranteapi/dishes/btcmarkets.go
// -------------------------------------
package btcmarkets

import (
	"fmt"
	"log"
	helper "restauranteapi/helper"
	"strconv"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// BalanceCrypto is to be exported
type BalanceCrypto struct {
	Balance        string // balance
	Currency       string // Currency
	CotacaoAtual   string // Cotacao
	ValueInCashAUD string // Value in AUD
	BestBid        string
	BestAsk        string
	LastPrice      string
	Instrument     string
	Volume24       string
	DateTime       string // date time
	Rotina         string // date time
	Date           string // date time
}

// CryptoCotacaoAdd is for export
func CryptoCotacaoAdd(redisclient *redis.Client, cryptoInsert BalanceCrypto) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "btccotacao"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	cryptoInsert.Date = cryptoInsert.DateTime[0:10]

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

// GetAll works
func GetAll(redisclient *redis.Client, currency string, rows string) []BalanceCrypto {

	rowsi, err := strconv.Atoi(rows)
	if rowsi == 0 {
		rowsi = 100
	}

	database := new(helper.DatabaseX)

	database.Collection = "btccotacao"

	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	// database.Database = "restaurante"
	// database.Location = "192.168.2.180"

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

	var results []BalanceCrypto

	// err1 := c.Find(bson.M{"name": dishName}).All(&result)
	// db.getCollection('btccotacao').find({}).sort({_id:-1}).limit(15);
	// err = c.Find(bson.M{"currency": currency}).All(&results)
	// err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
	// err = c.Find(bson.M{"currency": currency}).Sort("-_id").Limit(rowsi).All(&results)
	err = c.Find(bson.M{"currency": currency}).Sort("-datetime").Limit(rowsi).All(&results)
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

// GetDayStats works
func GetDayStats(redisclient *redis.Client, currency string, yearmonthday string, yearmonthdayend string) []BalanceCrypto {

	//"datetime" : "2018-01-01 12:12:40.877931086 +1100 AEDT m=+13.756967563"

	database := new(helper.DatabaseX)

	database.Collection = "btccotacao"

	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	// database.Database = "restaurante"
	// database.Location = "192.168.2.180"

	fmt.Println("database.Location")
	fmt.Println(database.Location)

	session, err := mgo.Dial(database.Location)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	var results []BalanceCrypto

	// db.users.findOne({"username" : {$regex : ".*son.*"}});
	// err = collection.Find(bson.M{  "date": yearmonthday,    "currency": currency}).All(&results)
	// err = collection.Find(bson.M{"datetime": bson.M{"$gt": yearmonthday}}).All(&results)

	//err = collection.Find(bson.M{"currency": currency, "datetime": bson.M{"$gte": yearmonthday, "$lte": yearmonthdayend}}).Sort("-_id").All(&results)
	err = collection.Find(bson.M{"currency": currency, "datetime": bson.M{"$gte": yearmonthday, "$lte": yearmonthdayend}}).Sort("-datetime").All(&results)

	// Below works in MongoDB Robo 3T interface
	//
	// db.getCollection('btccotacao').find({currency: {$in : ["ALL","XRP"]}})
	// db.getCollection('btccotacao').find({currency: {$eq : "ALL"}})
	// db.getCollection('btccotacao').find({datetime: /2018-01-07/ })
	// db.getCollection('btccotacao').find({datetime: {$gt : "2018-01-07"}})

	return results
}

// UpdateAllRows works
func UpdateAllRows(redisclient *redis.Client) []BalanceCrypto {

	// Not sure what to do ... perhaps just find a way to make the substring works :-)

	//"datetime" : "2018-01-01 12:12:40.877931086 +1100 AEDT m=+13.756967563"

	database := new(helper.DatabaseX)

	database.Collection = "btccotacao"

	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	// database.Database = "restaurante"
	// database.Location = "192.168.2.180"

	fmt.Println("database.Location")
	fmt.Println(database.Location)

	session, err := mgo.Dial(database.Location)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	var results []BalanceCrypto

	err = collection.Find(bson.M{}).All(&results)

	for x := 1; x < 10; x++ {

		err = collection.Find(bson.M{}).All(&results)
		results[x].Date = results[x].DateTime[0:10]
	}

	return results
}

// Import works
func Import(redisclient *redis.Client) []BalanceCrypto {

	databaseHome := new(helper.DatabaseX)

	databaseHome.Collection = "btccotacao"

	databaseHome.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	databaseHome.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	// Home PC
	databaseHome.Location = "192.168.2.170"

	session, err := mgo.Dial(databaseHome.Location)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(databaseHome.Database).C(databaseHome.Collection)

	var results []BalanceCrypto

	err = c.Find(nil).All(&results)
	if err != nil {
		// TODO: Do something about the error
		log.Fatal(err)
	} else {
		return results
	}

	return nil
}

// GetAllNoLimit brings it all
func GetAllNoLimit(redisclient *redis.Client, currency string) []BalanceCrypto {

	database := new(helper.DatabaseX)

	database.Collection = "btccotacao"

	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	// database.Database = "restaurante"
	// database.Location = "192.168.2.180"

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

	var results []BalanceCrypto

	err = c.Find(bson.M{"currency": currency}).Sort("-datetime").All(&results)
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
