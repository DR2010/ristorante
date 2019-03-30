// ...fjapidishes/fjapidishes.go
package main

import (
	"database/sql"
	"encoding/json"
	"fjapiactivities/helper"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

var mongodbvar helper.DatabaseX
var redisclient *redis.Client

var db *sql.DB
var err error
var sysid string

// Looks after the main routing
//
func main() {

	redisclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println(">>> Microservice: fjapiactivities running.")
	fmt.Println("Loading reference data in cache - Redis")

	sysid = helper.GetSYSID()

	loadreferencedatainredis()

	MSAPIactivitiesPort := helper.Getvaluefromcache("MSAPIactivitiesPort")
	MongoDBLocation := helper.Getvaluefromcache("API.MongoDB.Location")
	MongoDBDatabase := helper.Getvaluefromcache("API.MongoDB.Database")

	mongodbvar.Location = MongoDBLocation
	mongodbvar.Database = MongoDBDatabase

	fmt.Println("Running... Listening to " + MSAPIactivitiesPort)
	fmt.Println("MongoDB location: " + MongoDBLocation)
	fmt.Println("MongoDB database: " + MongoDBDatabase)

	router := XNewRouter()

	// handle using the router mux
	//
	http.Handle("/", router) // setting router rule

	err := http.ListenAndServe(":"+MSAPIactivitiesPort, nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", err)
	}
}

//#region Caching

// This is reading from ini file
//
func loadreferencedatainredis() {

	variable := helper.Readfileintostruct()

	err = redisclient.Set(sysid+"API.MongoDB.Location", variable.APIMongoDBLocation, 0).Err()
	err = redisclient.Set(sysid+"API.MongoDB.Database", variable.APIMongoDBDatabase, 0).Err()
	err = redisclient.Set(sysid+"API.APIServer.Port", variable.APIAPIServerPort, 0).Err()
	err = redisclient.Set(sysid+"API.APIServer.IPAddress", variable.APIAPIServerIPAddress, 0).Err()

	err = redisclient.Set(sysid+"Web.Debug", variable.WEBDebug, 0).Err()

	err = redisclient.Set(sysid+"CollectionOrders", variable.CollectionOrders, 0).Err()
	err = redisclient.Set(sysid+"CollectionSecurity", variable.CollectionSecurity, 0).Err()
	err = redisclient.Set(sysid+"CollectionDishes", variable.CollectionDishes, 0).Err()
	err = redisclient.Set(sysid+"CollectionActivities", variable.CollectionActivities, 0).Err()

	err = redisclient.Set(sysid+"MSAPIdishesPort", variable.MSAPIdishesPort, 0).Err()
	err = redisclient.Set(sysid+"MSAPIordersPort", variable.MSAPIordersPort, 0).Err()
	err = redisclient.Set(sysid+"MSAPIactivitiesPort", variable.MSAPIactivitiesPort, 0).Err()

}

type rediscachevalues struct {
	MongoDBLocation string
	MongoDBDatabase string
	APIServerPort   string
	APIServerIP     string
	WebDebug        string
}

//#endregion Caching

// Cache represents the cache data
type Cache struct {
	Key   string // cache key
	Value string // value in cache
}

func getcachedvalues(httpwriter http.ResponseWriter, req *http.Request) {

	var rv = new(rediscachevalues)

	rv.MongoDBLocation = helper.Getvaluefromcache("API.MongoDB.Location")
	rv.MongoDBDatabase = helper.Getvaluefromcache("API.MongoDB.Database")
	rv.APIServerPort = helper.Getvaluefromcache("API.APIServer.Port")
	rv.APIServerIP = helper.Getvaluefromcache("API.APIServer.IPAddress")
	rv.WebDebug = helper.Getvaluefromcache("Web.Debug")

	keys := make([]Cache, 5)
	keys[0].Key = "API.MongoDB.Location"
	keys[0].Value = rv.MongoDBLocation

	keys[1].Key = "API.MongoDB.Database"
	keys[1].Value = rv.MongoDBDatabase

	keys[2].Key = "API.APIServer.Port"
	keys[2].Value = rv.APIServerPort

	keys[3].Key = "API.APIServer.IPAddress"
	keys[3].Value = rv.APIServerIP

	keys[4].Key = "Web.Debug"
	keys[4].Value = rv.WebDebug

	json.NewEncoder(httpwriter).Encode(&keys)
}
