// Package dishes is a dish for packages
// -------------------------------------
// .../restauranteapi/activities/activities.go
// -------------------------------------
package activities

import (
	helper "fjapiactivities/helper"
	"fmt"
	"log"

	activities "fjapiactivities/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Add is for adding something (activities)
func Add(activityInsert activities.Activity) helper.Resultado {

	database := helper.GetDBParmFromCache("CollectionActivities")

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Insert(activityInsert)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Activity added"
	res.IsSuccessful = "Y"

	return res
}

// Find is to find stuff
func Find(objectFind string) (activities.Activity, string) {

	database := helper.GetDBParmFromCache("CollectionActivities")

	objectName := objectFind
	objectnull := activities.Activity{}

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	fmt.Println("Find objectname")
	fmt.Println(objectName)

	fmt.Println("Database")
	fmt.Println(database.Database)

	result := []activities.Activity{}
	err1 := c.Find(bson.M{"name": objectName}).All(&result)
	if err1 != nil {
		log.Fatal(err1)
	}

	var numrecsel = len(result)

	if numrecsel <= 0 {
		fmt.Println("404 Not found")

		return objectnull, "404 Not found"
	}

	fmt.Println("result[0]")
	fmt.Println(result[0])

	return result[0], "200 OK"
}

// Getall works
func Getall() []activities.Activity {

	database := helper.GetDBParmFromCache("CollectionActivities")

	session, err := mgo.Dial(database.Location)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	var results []activities.Activity

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

// GetAvailable works
func GetAvailable() []activities.Activity {

	database := helper.GetDBParmFromCache("CollectionActivities")

	session, err := mgo.Dial(database.Location)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	var results []activities.Activity

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

// Update is
func Update(objectUpdate activities.Activity) helper.Resultado {

	database := helper.GetDBParmFromCache("CollectionActivities")

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Update(bson.M{"name": objectUpdate.Name}, objectUpdate)

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "All good"
	res.IsSuccessful = "Y"

	if err != nil {

		fmt.Println(err)

		res.ErrorCode = "0002"
		res.ErrorDescription = err.Error()
		res.IsSuccessful = "Y"
	}

	return res
}

// Delete is
func Delete(objectDelete activities.Activity) helper.Resultado {

	database := helper.GetDBParmFromCache("CollectionActivities")

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Remove(bson.M{"name": objectDelete.Name})

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Activity deleted successfully"
	res.IsSuccessful = "Y"

	return res
}
