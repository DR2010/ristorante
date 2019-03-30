// Package main is the main package
// -------------------------------------
// .../fjapiactivities/activitieshandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	actions "fjapiactivities/activities"
	models "fjapiactivities/models"
	"fmt"
	"net/http"
)

// Hfind is
func Hfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	objectfound := models.Activity{}

	objecttofind := httprequest.FormValue("activityname") // This is the key, must be unique

	params := httprequest.URL.Query()
	parmactivityname := params.Get("activityname")

	fmt.Println("params.Get parmactivityname")
	fmt.Println(parmactivityname)

	fmt.Println("objecttofind")
	fmt.Println(objecttofind)

	objectfound, recordstatus := actions.Find(objecttofind)

	if recordstatus != "200 OK" {
		http.Error(httpwriter, "Record not found.", 400)
		return
	}

	json.NewEncoder(httpwriter).Encode(&objectfound)
}

// Hadd is
func Hadd(httpwriter http.ResponseWriter, req *http.Request) {

	toadd := models.Activity{}

	toadd.ShortName = req.FormValue("shortname")
	toadd.Description = req.FormValue("description") // This is the key, must be unique
	toadd.Type = req.FormValue("type")
	toadd.Status = req.FormValue("status")
	toadd.StartDate = req.FormValue("startdate")
	toadd.EndDate = req.FormValue("enddate")

	_, recordstatus := actions.Find(toadd.ShortName)
	if recordstatus == "200 OK" {
		http.Error(httpwriter, "Record already exists.", 422)
		return
	}

	// params := req.URL.Query()
	// dishtoadd.Name = params.Get("dishname")
	// dishtoadd.Type = params.Get("dishtype")
	// dishtoadd.Price = params.Get("dishprice")
	// dishtoadd.GlutenFree = params.Get("dishglutenfree")
	// dishtoadd.DairyFree = params.Get("dishdairyfree")
	// dishtoadd.Vegetarian = params.Get("dishvegetarian")

	ret := actions.Add(toadd)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hupdate is
func Hupdate(httpwriter http.ResponseWriter, req *http.Request) {

	objectaction := getobject(httpwriter, req)

	ret := actions.Update(objectaction)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hdelete is
func Hdelete(httpwriter http.ResponseWriter, req *http.Request) {

	objectaction := getobject(httpwriter, req)

	ret := actions.Delete(objectaction)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// getobject is
func getobject(httpwriter http.ResponseWriter, req *http.Request) models.Activity {

	objectaction := models.Activity{}

	objectaction.ShortName = req.FormValue("activityshorthname") // This is the key, must be unique
	objectaction.Type = req.FormValue("activitytype")
	objectaction.Description = req.FormValue("activitydescription")
	objectaction.StartDate = req.FormValue("activitystartdate")
	objectaction.EndDate = req.FormValue("activityenddate")
	objectaction.LongDescription = req.FormValue("activitylongdescription")
	objectaction.Status = req.FormValue("activitystatus")

	return objectaction
}

// Halsolist is
func Halsolist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = actions.Getall()

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

// Hlist is a function to return a list of dishes
func Hlist(httpwriter http.ResponseWriter, req *http.Request) {

	var list = actions.Getall()

	json.NewEncoder(httpwriter).Encode(&list)
}

// Hlistavailable is a function to return a list of dishes
func Hlistavailable(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = actions.GetAvailable()

	json.NewEncoder(httpwriter).Encode(&dishlist)
}
