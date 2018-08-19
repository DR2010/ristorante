// Package main is the main package
// -------------------------------------
// .../restauranteapi/orderhandler.go
// -------------------------------------
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"restauranteapi/helper"
	orders "restauranteapi/orders"
	"restauranteapi/security"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Hfind finds orders
func Hfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	objfound := orders.Order{}

	objtofind := httprequest.FormValue("orderid") // This is the key, must be unique

	objfound, _ = orders.Find(redisclient, objtofind)

	json.NewEncoder(httpwriter).Encode(&objfound)
}

// Horderfind is
func Horderfind(httpwriter http.ResponseWriter, httprequest *http.Request) {

	redisclient := helper.GetRedisPointer()

	orderfound := orders.Order{}

	ordertofind := httprequest.FormValue("orderid") // This is the key, must be unique

	params := httprequest.URL.Query()
	parmorderid := params.Get("orderid")

	fmt.Println("params.Get parmorderid")
	fmt.Println(parmorderid)

	orderfound, _ = orders.Find(redisclient, ordertofind)

	json.NewEncoder(httpwriter).Encode(&orderfound)
}

// Horderadd add orders
func Horderadd(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)
	// bodystr := string(bodybyte[:])

	type dcOrderItem struct {
		Pratoname  string //
		Quantidade string //
		Price      string //
		Total      string //
	}

	type dcOrder struct {
		ID         string // random ID for order, yet to define algorithm
		ClientName string // Client Name for the order
		ClientID   string // Client ID in case they logon - later
		Date       string // Order Date
		Time       string // Order Time
		EatMode    string // Delivery, Eat In, Take Away
		Status     string // Status
		Items      []dcOrderItem
	}

	var objtoaction dcOrder
	err = json.Unmarshal(bodybyte, &objtoaction)

	tries := 1
	for tries < 1000 {

		rand.Seed(time.Now().UTC().UnixNano())
		objtoaction.ID = strconv.Itoa(rand.Intn(100000))

		_, recordstatus := orders.Find(redisclient, objtoaction.ID)

		if recordstatus == "200 OK" {
			fmt.Println("recordstatus")
			fmt.Println(recordstatus)
			http.Error(httpwriter, "Record already exists.", 422)
			fmt.Println("try=" + strconv.Itoa(tries))
			tries++
			continue
		}
		break
	}

	objtoactionMAP := orders.Order{}
	objtoactionMAP.ID = objtoaction.ID
	objtoactionMAP.ClientID = objtoaction.ClientID
	objtoactionMAP.ClientName = objtoaction.ClientName
	objtoactionMAP.Date = objtoaction.Date
	objtoactionMAP.Time = objtoaction.Time
	objtoactionMAP.Status = objtoaction.Status
	objtoactionMAP.EatMode = objtoaction.EatMode

	var slen = len(objtoaction.Items)
	objtoactionMAP.Items = make([]orders.Item, slen)

	var totalgeral = 0.00

	// I have to remove the header coming from the caller.
	// Perhaps the caller should suppress the header somehow

	var destindex = 0

	for index, element := range objtoaction.Items {
		// index is the index where we are
		// element is the element from someSlice for where we are

		// if index == 0 {
		// 	continue
		// }

		// destindex = index - 1

		destindex = index

		objtoactionMAP.Items[destindex].PratoName = element.Pratoname
		objtoactionMAP.Items[destindex].Price = element.Price
		objtoactionMAP.Items[destindex].Quantidade = element.Quantidade
		objtoactionMAP.Items[destindex].Total = element.Total

		prc, _ := strconv.ParseFloat(element.Price, 64)
		qty, _ := strconv.ParseFloat(element.Quantidade, 64)
		tot := prc * qty
		totalgeral = totalgeral + tot
		// objtoactionMAP.Items[destindex].Total = strconv.Itoa(tot)

	}
	// objtoactionMAP.TotalGeral = strconv.Itoa(totalgeral)
	// objtoactionMAP.TotalGeral = strconv.FormatFloat(totalgeral, 'g', -1, 64)

	objtoactionMAP.TotalGeral = fmt.Sprintf("%.2f", totalgeral)

	ret := orders.Add(redisclient, objtoactionMAP)

	if ret.IsSuccessful == "Y" {
		// do something

		fmt.Println("Order added successfully:" + objtoaction.ClientName)

		type RespAddOrder struct {
			ID string
		}

		// return value
		obj := &RespAddOrder{ID: objtoaction.ID}
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpwriter, string(bresp)) // write data to response
	}

	return
}

// HAPIorderadd add orders
// This is the V2 to handle a generation of User ID
func HAPIorderadd(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)
	// bodystr := string(bodybyte[:])

	type dcOrderItem struct {
		Pratoname  string //
		Quantidade string //
		Price      string //
		Total      string //
	}

	type dcOrder struct {
		ID         string // random ID for order, yet to define algorithm
		ClientName string // Client Name for the order
		ClientID   string // Client ID in case they logon - later
		Date       string // Order Date
		Time       string // Order Time
		EatMode    string // Delivery, Eat In, Take Away
		Status     string // Status
		Items      []dcOrderItem
	}

	var objtoaction dcOrder
	err = json.Unmarshal(bodybyte, &objtoaction)

	// Generate order ID
	objtoaction.ID = GenerateNewOrderID()

	// If user ID is passed in, just use it
	if objtoaction.ClientID == "Anonymous" || objtoaction.ClientID == "" {

		// Create the user
		//
		credentials := security.Credentials{}
		credentials.UserID = objtoaction.ClientID
		credentials.ApplicationID = "Restaurante"
		credentials.Name = objtoaction.ClientName
		credentials.Password = "NA"
		credentials.IsAdmin = "N"

		// Generate user ID
		objtoaction.ClientID = GenerateNewUser()
		credentials.UserID = objtoaction.ClientID

		security.Useradd(redisclient, credentials)

		// Generate User ID if it is set to Anonymous
		// USR+99999
		// ------------------------------------------
	}

	objtoactionMAP := orders.Order{}
	objtoactionMAP.ID = objtoaction.ID
	objtoactionMAP.ClientID = objtoaction.ClientID
	objtoactionMAP.ClientName = objtoaction.ClientName
	objtoactionMAP.Date = objtoaction.Date
	objtoactionMAP.Time = objtoaction.Time
	objtoactionMAP.Status = objtoaction.Status
	objtoactionMAP.EatMode = objtoaction.EatMode

	var slen = len(objtoaction.Items)
	objtoactionMAP.Items = make([]orders.Item, slen)

	var totalgeral = 0.00

	// I have to remove the header coming from the caller.
	// Perhaps the caller should suppress the header somehow

	var destindex = 0

	for index, element := range objtoaction.Items {
		// index is the index where we are
		// element is the element from someSlice for where we are

		// if index == 0 {
		// 	continue
		// }

		// destindex = index - 1

		destindex = index

		objtoactionMAP.Items[destindex].PratoName = element.Pratoname
		objtoactionMAP.Items[destindex].Price = element.Price
		objtoactionMAP.Items[destindex].Quantidade = element.Quantidade
		objtoactionMAP.Items[destindex].Total = element.Total

		prc, _ := strconv.ParseFloat(element.Price, 64)
		qty, _ := strconv.ParseFloat(element.Quantidade, 64)
		tot := prc * qty
		totalgeral = totalgeral + tot
		// objtoactionMAP.Items[destindex].Total = strconv.Itoa(tot)

	}
	// objtoactionMAP.TotalGeral = strconv.Itoa(totalgeral)
	// objtoactionMAP.TotalGeral = strconv.FormatFloat(totalgeral, 'g', -1, 64)

	objtoactionMAP.TotalGeral = fmt.Sprintf("%.2f", totalgeral)

	ret := orders.Add(redisclient, objtoactionMAP)

	if ret.IsSuccessful == "Y" {
		// do something

		fmt.Println("Order added successfully:" + objtoaction.ClientName)

		type RespAddOrder struct {
			ID       string
			ClientID string
		}

		objret := RespAddOrder{}
		objret.ID = objtoaction.ID
		objret.ClientID = objtoaction.ClientID

		// return value
		// obj := &RespAddOrder{ID: objtoaction.ID, ClientID: objtoaction.ClientID}
		// obj := &RespAddOrder{ID: objtoaction.ID}

		bresp, _ := json.Marshal(objret)

		fmt.Fprintf(httpwriter, string(bresp)) // write data to response
	}

	return
}

// GenerateNewOrderID is to generate a new order
// ----------------------------------------------
func GenerateNewOrderID() string {

	orderid := "0"

	// Generate order ID
	// ------------------
	tries := 1
	for tries < 1000 {

		rand.Seed(time.Now().UTC().UnixNano())
		orderid = strconv.Itoa(rand.Intn(100000))

		_, recordstatus := orders.Find(redisclient, orderid)

		if recordstatus == "200 OK" {
			fmt.Println("recordstatus")
			fmt.Println(recordstatus)
			fmt.Println("try=" + strconv.Itoa(tries))
			tries++
			continue
		}
		break
	}

	return orderid
}

// GenerateNewUser is to generate a new order
// ----------------------------------------------
func GenerateNewUser() string {

	userid := ""

	// Generate User ID
	// ------------------
	tries := 1
	for tries < 1000 {

		rand.Seed(time.Now().UTC().UnixNano())
		userid = strconv.Itoa(rand.Intn(100000))

		_, recordstatus := security.Find(redisclient, userid)

		if recordstatus == "200 OK" {
			fmt.Println("recordstatus")
			fmt.Println(recordstatus)
			fmt.Println("try=" + strconv.Itoa(tries))
			tries++
			continue
		}

		break
	}

	return "USR" + userid
}

// Horderupdate update orders
func Horderupdate(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)
	// bodystr := string(bodybyte[:])

	type dcOrderItem struct {
		Pratoname  string //
		Quantidade string //
		Preco      string //
	}

	// Esta estrutura e' usada pelo Javascript para adicionar e chamar a API
	// Tem que manter a estrutura do Javascript in sync com o golang
	// Agora nao esta em sync. 8/2/2018

	// type dcOrder struct {
	// 	OrderID         string // random ID for order, yet to define algorithm
	// 	OrderClientID   string // Client ID in case they logon - later
	// 	OrderClientName string // Client Name for the order
	// 	OrderDate       string // Order Date
	// 	OrderTime       string // Order Time
	// 	EatMode         string // Delivery, Eat In, Take Away
	// 	Status          string // Status
	// 	Pratos          []dcOrderItem
	// }

	var objtoaction orders.Order
	err = json.Unmarshal(bodybyte, &objtoaction)

	_, recordstatus := orders.Find(redisclient, objtoaction.ID)

	if recordstatus == "200 OK" {
		fmt.Println("recordstatus")
		fmt.Println(recordstatus)
	}

	objtoactionMAP := orders.Order{}
	objtoactionMAP.ID = objtoaction.ID
	objtoactionMAP.ClientID = objtoaction.ClientID
	objtoactionMAP.ClientName = objtoaction.ClientName
	objtoactionMAP.Date = objtoaction.Date
	objtoactionMAP.Time = objtoaction.Time
	objtoactionMAP.Status = objtoaction.Status
	objtoactionMAP.EatMode = objtoaction.EatMode

	var slen = len(objtoaction.Items)
	objtoactionMAP.Items = make([]orders.Item, slen)

	var totalgeral = 0.00

	// I have to remove the header coming from the caller.
	// Perhaps the caller should suppress the header somehow

	var destindex = 0

	for index, element := range objtoaction.Items {
		// index is the index where we are
		// element is the element from someSlice for where we are

		// if index == 0 {
		// 	continue
		// }

		// destindex = index - 1
		// destindex = index

		// objtoactionMAP.Items[destindex].PratoName = element.PratoName
		// objtoactionMAP.Items[destindex].Price = element.Price
		// objtoactionMAP.Items[destindex].Quantidade = element.Quantidade

		// prc, _ := strconv.Atoi(element.Price)
		// qty, _ := strconv.Atoi(element.Price)
		// tot := prc * qty
		// totalgeral = totalgeral + tot

		// objtoactionMAP.Items[destindex].Total = strconv.Itoa(tot)

		destindex = index

		objtoactionMAP.Items[destindex].PratoName = element.PratoName
		objtoactionMAP.Items[destindex].Price = element.Price
		objtoactionMAP.Items[destindex].Quantidade = element.Quantidade
		objtoactionMAP.Items[destindex].Total = element.Total

		prc, _ := strconv.ParseFloat(element.Price, 64)
		qty, _ := strconv.ParseFloat(element.Quantidade, 64)
		tot := prc * qty
		totalgeral = totalgeral + tot

	}
	// objtoactionMAP.TotalGeral = strconv.Itoa(totalgeral)
	objtoactionMAP.TotalGeral = fmt.Sprintf("%.2f", totalgeral)

	ret := orders.Update(redisclient, objtoactionMAP)

	if ret.IsSuccessful == "Y" {
		// do something

		fmt.Println("Order added successfully:" + objtoaction.ClientName)

		type RespAddOrder struct {
			ID string
		}

		// return value
		obj := &RespAddOrder{ID: objtoaction.ID}
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpwriter, string(bresp)) // write data to response
	}

	return
}

// ------------------------------------------------------------------
// ------------------------------------------------------------------
// ------------------------------------------------------------------
// ------------------------------------------------------------------
// ------------------------------------------------------------------
// CRIAR um UPDATE apenas para o status
// porem tenho que consertar o UPDATE ALL FIELDS pois esta removendo os TOTAIS !!!!!!!!!!!!!!!!!!!! 11/02/2018
// ------------------------------------------------------------------
// ------------------------------------------------------------------
// ------------------------------------------------------------------
// ------------------------------------------------------------------
// ------------------------------------------------------------------

// Hdelete delete orders
func Hdelete(httpwriter http.ResponseWriter, req *http.Request) {

	objtoupdate := orders.Order{}

	objtoupdate.ClientID = req.FormValue("orderID") // This is the key, must be unique

	ret := orders.Delete(redisclient, objtoupdate.ClientID)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Halsolist list orders
func Halsolist(httpwriter http.ResponseWriter, req *http.Request) {

	var orderlist = orders.Getall(redisclient)

	json.NewEncoder(httpwriter).Encode(&orderlist)
}

// OrderList also list orders
func OrderList(httpwriter http.ResponseWriter, req *http.Request) {

	var orderlist = orders.Getall(redisclient)
	json.NewEncoder(httpwriter).Encode(&orderlist)
}

// OrderListV2 also list orders
func OrderListV2(httpwriter http.ResponseWriter, req *http.Request) {

	var userid = req.FormValue("clientid") // This is the key, must be unique

	if userid == "" {
		// orderlist1 := orders.Getall(redisclient)
		orderlist1 := orders.Getallbutcompleted(redisclient)
		json.NewEncoder(httpwriter).Encode(&orderlist1)
	} else {
		orderlist2 := orders.GetallbyUser(redisclient, userid)
		json.NewEncoder(httpwriter).Encode(&orderlist2)
	}

	// json.NewEncoder(httpwriter).Encode(&orderlist)
}

// CopyOrdersToMySQL also list orders
func CopyOrdersToMySQL(httpwriter http.ResponseWriter, req *http.Request) {

	db, err = sql.Open("mysql", "daniel:oculos18@/festajunina")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	orders.SavetoMySQL(redisclient, db)
}

// ordercompleted also list orders
func ordercompleted(httpwriter http.ResponseWriter, req *http.Request) {

	status := "Completed"
	orderlist2 := orders.Getallcompleted(redisclient, status)
	json.NewEncoder(httpwriter).Encode(&orderlist2)

	// json.NewEncoder(httpwriter).Encode(&orderlist)
}

// orderstatus also list orders
func orderstatus(httpwriter http.ResponseWriter, req *http.Request) {

	var status = req.FormValue("status")

	orderlist2 := orders.Getallcompleted(redisclient, status)
	json.NewEncoder(httpwriter).Encode(&orderlist2)

	// json.NewEncoder(httpwriter).Encode(&orderlist)
}
