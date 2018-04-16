// Package main is the main package
// -------------------------------------
// .../restauranteapi/dishhandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	dishesmethods "restauranteapi/dishes"
	helper "restauranteapi/helper"
	dishes "restauranteapi/models"
)

// Hdishfind is
func Hdishfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	redisclient := helper.GetRedisPointer()

	dishfound := dishes.Dish{}

	dishtofind := httprequest.FormValue("dishname") // This is the key, must be unique

	params := httprequest.URL.Query()
	parmdishname := params.Get("dishname")

	fmt.Println("params.Get parmdishname")
	fmt.Println(parmdishname)

	fmt.Println("httprequest.FormValue dishname")
	fmt.Println(dishtofind)

	dishfound, _ = dishesmethods.Find(redisclient, dishtofind)

	json.NewEncoder(httpwriter).Encode(&dishfound)
}

// Hdishadd is
func Hdishadd(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")
	dishtoadd.InitialAvailable = req.FormValue("dishinitialavailable")
	dishtoadd.CurrentAvailable = req.FormValue("dishcurrentavailable")
	dishtoadd.ImageName = req.FormValue("dishimagename")
	dishtoadd.Description = req.FormValue("dishdescription")
	dishtoadd.Descricao = req.FormValue("dishdescricao")

	_, recordstatus := dishesmethods.Find(redisclient, dishtoadd.Name)
	if recordstatus == "200 OK" {
		fmt.Println("dishtoadd.Name")
		fmt.Println(dishtoadd.Name)

		fmt.Println("recordstatus")
		fmt.Println(recordstatus)
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

	ret := dishesmethods.Dishadd(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hdishupdate is
func Hdishupdate(httpwriter http.ResponseWriter, req *http.Request) {

	redisclient := helper.GetRedisPointer()

	dishtoupdate := dishes.Dish{}

	dishtoupdate.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoupdate.Type = req.FormValue("dishtype")
	dishtoupdate.Price = req.FormValue("dishprice")
	dishtoupdate.GlutenFree = req.FormValue("dishglutenfree")
	dishtoupdate.DairyFree = req.FormValue("dishdairyfree")
	dishtoupdate.Vegetarian = req.FormValue("dishvegetarian")
	dishtoupdate.InitialAvailable = req.FormValue("dishinitialavailable")
	dishtoupdate.CurrentAvailable = req.FormValue("dishcurrentavailable")
	dishtoupdate.ImageName = req.FormValue("dishimagename")
	dishtoupdate.Description = req.FormValue("dishdescription")
	dishtoupdate.Descricao = req.FormValue("dishdescricao")

	fmt.Println("dishtoupdate.Name")
	fmt.Println(dishtoupdate.Name)

	// params := req.URL.Query()
	// dishtoadd.Name = params.Get("dishname")
	// dishtoadd.Type = params.Get("dishtype")
	// dishtoadd.Price = params.Get("dishprice")
	// dishtoadd.GlutenFree = params.Get("dishglutenfree")
	// dishtoadd.DairyFree = params.Get("dishdairyfree")
	// dishtoadd.Vegetarian = params.Get("dishvegetarian")

	ret := dishesmethods.Dishupdate(redisclient, dishtoupdate)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hdishdelete is
func Hdishdelete(httpwriter http.ResponseWriter, req *http.Request) {

	redisclient := helper.GetRedisPointer()

	dishtoupdate := dishes.Dish{}

	dishtoupdate.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoupdate.Type = req.FormValue("dishtype")
	dishtoupdate.Price = req.FormValue("dishprice")
	dishtoupdate.GlutenFree = req.FormValue("dishglutenfree")
	dishtoupdate.DairyFree = req.FormValue("dishdairyfree")
	dishtoupdate.Vegetarian = req.FormValue("dishvegetarian")
	dishtoupdate.InitialAvailable = req.FormValue("dishinitialavailable")
	dishtoupdate.CurrentAvailable = req.FormValue("dishcurrentavailable")
	dishtoupdate.ImageName = req.FormValue("dishimagename")
	dishtoupdate.Description = req.FormValue("dishdescription")
	dishtoupdate.Descricao = req.FormValue("dishdescricao")

	ret := dishesmethods.Dishdelete(redisclient, dishtoupdate)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hdishalsolist is
func Hdishalsolist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishesmethods.Getall(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

// Hdishlist is a function to return a list of dishes
func Hdishlist(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishesmethods.Getall(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}

// Hdishlistavailable is a function to return a list of dishes
func Hdishlistavailable(httpwriter http.ResponseWriter, req *http.Request) {

	var dishlist = dishesmethods.GetAvailable(redisclient)

	json.NewEncoder(httpwriter).Encode(&dishlist)
}
