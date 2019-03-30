package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis"
)

var redisclient *redis.Client
var SYSID string
var databaseEV DatabaseX

// DatabaseX is a struct
type DatabaseX struct {
	Location   string // location of the database localhost, something.com, etc
	Database   string // database name
	Collection string // collection name
}

// Resultado is a struct
type Resultado struct {
	ErrorCode        string // error code
	ErrorDescription string // description
	IsSuccessful     string // Y or N
	ReturnedValue    string // Any string
}

// GetRedisPointer returns
func GetRedisPointer(bucket int) *redis.Client {

	bucket = 0

	if redisclient == nil {
		redisclient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",     // no password set
			DB:       bucket, // use default DB
		})
	}

	return redisclient
}

// RestEnvVariables = restaurante environment variables
//
type RestEnvVariables struct {
	APIMongoDBLocation    string // location of the database localhost, something.com, etc
	APIMongoDBDatabase    string // database name
	APIAPIServerPort      string // collection name
	APIAPIServerIPAddress string // apiserver name
	WEBDebug              string // debug
	CollectionOrders      string // Collection Names
	CollectionSecurity    string // Collection Names
	CollectionDishes      string // Collection Names
	CollectionActivities  string // Collection Names
	MSAPIdishesPort       string // Microservices Port Dishes
	MSAPIordersPort       string // Microservices Port Orders
	MSAPIactivitiesPort   string // Microservices Port Orders
	MSAPItemperaturePort  string // Microservices Port temperature
	SYSID                 string // Collection Names

}

// Readfileintostruct is
func Readfileintostruct() RestEnvVariables {
	dat, err := ioutil.ReadFile("fjapiactivities.ini")
	check(err)
	fmt.Print(string(dat))

	var restenv RestEnvVariables

	json.Unmarshal(dat, &restenv)

	return restenv
}

// GetSYSID is just returning the System ID directly from file
// It is happening to enable multiple usage of Redis Keys ("SYSID" + "APIURL" for instance)
func GetSYSID() string {

	if SYSID == "" {

		dat, err := ioutil.ReadFile("fjapiactivities.ini")
		check(err)
		// fmt.Print(string(dat))

		var restenv RestEnvVariables

		json.Unmarshal(dat, &restenv)

		SYSID = restenv.SYSID

		return restenv.SYSID
	}

	return SYSID

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Getvaluefromcache returns the value of a key from cache
func Getvaluefromcache(key string) string {

	// bucket is ZERO for now
	// I am allowing it to be setup now
	rp := GetRedisPointer(0)

	sysid := GetSYSID()

	valuetoreturn, _ := rp.Get(sysid + key).Result()

	return valuetoreturn
}

// GetDBParmFromCache returns the value of a key from cache
func GetDBParmFromCache(collection string) *DatabaseX {

	database := new(DatabaseX)

	database.Collection = Getvaluefromcache(collection)
	database.Database = Getvaluefromcache("API.MongoDB.Database")
	database.Location = Getvaluefromcache("API.MongoDB.Location")

	return database
}
