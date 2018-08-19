// Package business is a business for packages
// -------------------------------------
// .../restauranteapi/business/borganisation.go
// -------------------------------------
package business

import (
	"fmt"
	"log"
	helper "restauranteapi/helper"

	models "restauranteapi/models"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OrganisationAdd is for export
func OrganisationAdd(redisclient *redis.Client, objInsert models.Organisation) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "organisations"
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

	err = collection.Insert(objInsert)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Event added"
	res.IsSuccessful = "Y"

	return res
}

// OrganisationFind is to find stuff
func OrganisationFind(redisclient *redis.Client, objFind string) (models.Organisation, string) {

	database := new(helper.DatabaseX)
	database.Collection = "organisations"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	eventID := objFind
	eventnull := models.Organisation{}

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	result := []models.Organisation{}
	err1 := c.Find(bson.M{"ID": eventID}).All(&result)
	if err1 != nil {
		log.Fatal(err1)
	}

	var numrecsel = len(result)

	if numrecsel <= 0 {
		return eventnull, "404 Not found"
	}

	return result[0], "200 OK"
}

// OrganisationGetall works
func OrganisationGetall(redisclient *redis.Client) []models.Organisation {

	database := new(helper.DatabaseX)

	database.Collection = "organisations"

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

	var results []models.Organisation

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

// OrganisationGetAvailable works
func OrganisationGetAvailable(redisclient *redis.Client) []models.Organisation {

	database := new(helper.DatabaseX)

	database.Collection = "organisations"

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

	var results []models.Organisation

	err = c.Find(bson.M{"currentavailable": bson.M{"$ne": "0"}}).All(&results)

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

// OrganisationUpdate is
func OrganisationUpdate(redisclient *redis.Client, organisationUpdate models.Organisation) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "organisations"
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

	err = collection.Update(bson.M{"ID": organisationUpdate.ID}, organisationUpdate)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Something Happened"
	res.IsSuccessful = "Y"

	return res
}

// OrganisationDelete is
func OrganisationDelete(redisclient *redis.Client, organisationDelete models.Organisation) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "organisations"
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

	err = collection.Remove(bson.M{"ID": organisationDelete.ID})

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Organisation deleted successfully"
	res.IsSuccessful = "Y"

	return res
}
