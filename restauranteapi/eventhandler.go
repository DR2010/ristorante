// Package main is the main package
// -------------------------------------
// .../restauranteapi/eventhandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	eventsmethods "restauranteapi/business"
	dishesmethods "restauranteapi/dishes"
	helper "restauranteapi/helper"
	events "restauranteapi/models"
)

// Heventfind is
func Heventfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	redisclient := helper.GetRedisPointer()

	eventfound := events.Event{}

	eventtofind := httprequest.FormValue("eventname") // This is the key, must be unique

	params := httprequest.URL.Query()
	parmeventname := params.Get("eventname")

	fmt.Println("params.Get parmeventname")
	fmt.Println(parmeventname)

	fmt.Println("httprequest.FormValue eventname")
	fmt.Println(eventtofind)

	eventfound, _ = eventsmethods.EventFind(redisclient, eventtofind)

	json.NewEncoder(httpwriter).Encode(&eventfound)
}

// Heventadd is
func Heventadd(httpwriter http.ResponseWriter, req *http.Request) {

	eventtoadd := events.Event{}

	eventtoadd.ID = req.FormValue("eventid") // This is the key, must be unique
	eventtoadd.Name = req.FormValue("eventname")
	eventtoadd.Description = req.FormValue("eventdescription")
	eventtoadd.Location = req.FormValue("eventlocation")
	eventtoadd.Date = req.FormValue("eventdate")
	eventtoadd.Manager = req.FormValue("eventmanager")
	eventtoadd.Mobile = req.FormValue("eventmobile")

	_, recordstatus := dishesmethods.Find(redisclient, eventtoadd.Name)
	if recordstatus == "200 OK" {
		fmt.Println("eventtoadd.Name")
		fmt.Println(eventtoadd.Name)

		fmt.Println("recordstatus")
		fmt.Println(recordstatus)
		http.Error(httpwriter, "Event alread registered.", 422)
		return
	}

	// params := req.URL.Query()
	// eventtoadd.ID = params.Get("eventid")
	// eventtoadd.Name = params.Get("eventname")
	// eventtoadd.Description = params.Get("eventdescription")
	// eventtoadd.Location = params.Get("eventlocation")
	// eventtoadd.Date = params.Get("eventdate")

	ret := eventsmethods.EventAdd(redisclient, eventtoadd)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Heventupdate is
func Heventupdate(httpwriter http.ResponseWriter, req *http.Request) {

	redisclient := helper.GetRedisPointer()

	eventtoupdate := events.Event{}

	eventtoupdate.ID = req.FormValue("eventid") // This is the key, must be unique
	eventtoupdate.Name = req.FormValue("eventname")
	eventtoupdate.Description = req.FormValue("eventdescription")
	eventtoupdate.Location = req.FormValue("eventlocation")
	eventtoupdate.Date = req.FormValue("eventdate")

	fmt.Println("eventtoupdate.Name")
	fmt.Println(eventtoupdate.Name)

	ret := eventsmethods.EventUpdate(redisclient, eventtoupdate)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Heventdelete is
func Heventdelete(httpwriter http.ResponseWriter, req *http.Request) {

	redisclient := helper.GetRedisPointer()

	eventtoupdate := events.Event{}

	eventtoupdate.ID = req.FormValue("eventid") // This is the key, must be unique
	eventtoupdate.Name = req.FormValue("eventname")
	eventtoupdate.Description = req.FormValue("eventdescription")
	eventtoupdate.Location = req.FormValue("eventlocation")
	eventtoupdate.Date = req.FormValue("eventdate")

	ret := eventsmethods.EventDelete(redisclient, eventtoupdate)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Heventalsolist is
func Heventalsolist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishesmethods.Getall(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

// Heventlist is a function to return a list of dishes
func Heventlist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = eventsmethods.Getall(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

// Heventlistavailable is a function to return a list of dishes
func Heventlistavailable(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishesmethods.GetAvailable(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}
