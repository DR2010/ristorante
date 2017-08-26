package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"mongodb/dishes"
	"mongodb/helper"
	"mongodb/order"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

var mongodbvar helper.DatabaseX

var db *sql.DB
var err error

// Looks after the main routing
//
func main() {

	// db, err = sql.Open("mysql", "root:oculos18@/gufcdraws")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err.Error())
	// }

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"

	fmt.Println("Running... Listening to :1515")

	router := XNewRouter()

	// handle using the router mux
	//
	http.Handle("/", router) // setting router rule

	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	http.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./fonts"))))

	err := http.ListenAndServe(":1515", nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", router)
	}
}

func root(httpwriter http.ResponseWriter, r *http.Request) {

	// create new template
	var listtemplate = `
		{{define "listtemplate"}}
	
		{{end}}
		`

	t, _ := template.ParseFiles("templates/indextemplate.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpwriter, listtemplate)
	return
}

func root2(httpwriter http.ResponseWriter, r *http.Request) {
	http.ServeFile(httpwriter, r, "index.html")

	return
}

// How to access parms from the URL
func printparm(httpwriter http.ResponseWriter, r *http.Request) {

	// extracts parm values from URL
	//
	r.ParseForm()

	//Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	// fmt.Fprintf(w, "Hello astaxie!") // write data to response

	// httpwriter.Write([]byte("nothing\n"))
	// httpwriter.Write([]byte("path" + r.URL.Path))
	// httpwriter.Write([]byte("scheme" + r.URL.Scheme))

	for key, v := range r.Form {
		httpwriter.Write([]byte("\n" + key))
		httpwriter.Write([]byte(" : " + strings.Join(v, "")))
	}

}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")

		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/loginPage", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/loginPage", 301)
		return
	}

	res.Write([]byte("Hello" + databaseUsername))

}

func dishlist(httpwriter http.ResponseWriter, req *http.Request) {

	type ControllerInfo struct {
		Name string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/listtemplate.html")

	// var dishlist = []dishes.Dish{
	// 	dishes.Dish{Type: "Main", Name: "Feijoada", Price: "200.00"},
	// 	dishes.Dish{Type: "Main", Name: "Batatada", Price: "130.00"},
	// 	dishes.Dish{Type: "Desert", Name: "Pudim de leite", Price: "50.00"},
	// 	dishes.Dish{Type: "Entree", Name: "Almonds", Price: "25.00"},
	// }

	var dishlist = dishes.GetAll(mongodbvar)

	items := DisplayTemplate{}
	items.Info.Name = "Dish List"

	var numberoffields = 6

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Name"
	items.FieldNames[1] = "Type"
	items.FieldNames[2] = "Price"
	items.FieldNames[3] = "GlutenFree"
	items.FieldNames[4] = "DairyFree"
	items.FieldNames[5] = "Vegetarian"

	// Set rows to be displayed
	items.Rows = make([]Row, len(dishlist))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(dishlist); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = dishlist[i].Name
		items.Rows[i].Description[1] = dishlist[i].Type
		items.Rows[i].Description[2] = dishlist[i].Price
		items.Rows[i].Description[3] = dishlist[i].GlutenFree
		items.Rows[i].Description[4] = dishlist[i].DairyFree
		items.Rows[i].Description[5] = dishlist[i].Vegetarian
	}

	t.Execute(httpwriter, items)
}

func orderlist(httpwriter http.ResponseWriter, req *http.Request) {

	type ControllerInfo struct {
		Name string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/listtemplate.html")

	var orderlist = []order.Order{
		order.Order{ID: "Main", ClientName: "Daniel", Date: "10-Jan-2017", Time: "10:00", Total: "100.00"},
		order.Order{ID: "Main", ClientName: "Katia", Date: "10-Jan-2017", Time: "10:00", Total: "100.00"},
		order.Order{ID: "Main", ClientName: "Arthur", Date: "10-Jan-2017", Time: "10:00", Total: "100.00"},
		order.Order{ID: "Main", ClientName: "Kevin", Date: "10-Jan-2017", Time: "10:00", Total: "100.00"},
	}

	numitems := len(orderlist)
	numcolumns := 5

	items := DisplayTemplate{}
	items.Info.Name = "Order List"

	// Set colum names
	items.FieldNames = make([]string, numcolumns)
	items.FieldNames[0] = "ID"
	items.FieldNames[1] = "ClientName"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Time"
	items.FieldNames[4] = "Total"

	// Set rows to be displayed
	items.Rows = make([]Row, numitems)

	for i := 0; i < numitems; i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numcolumns)
		items.Rows[i].Description[0] = orderlist[i].ID
		items.Rows[i].Description[1] = orderlist[i].ClientName
		items.Rows[i].Description[2] = orderlist[i].Date
		items.Rows[i].Description[3] = orderlist[i].Time
		items.Rows[i].Description[4] = orderlist[i].Total
	}

	t.Execute(httpwriter, items)

}

func dishadddisplay(httpwriter http.ResponseWriter, req *http.Request) {

	type ControllerInfo struct {
		Name string
	}

	type Row struct {
		Description []string
	}

	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishadd.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	t.Execute(httpwriter, items)
	return

}

func dishadd(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := dishes.Dishadd(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

func dishupdatedisplay(httpwriter http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	// Get all selected records
	dishselected := req.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	type ControllerInfo struct {
		Name string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		DishItem   dishes.Dish
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishupdate.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	items.DishItem = dishes.Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = dishes.Dish{}
	var dishname = items.DishItem.Name

	dishfind = dishes.Find(mongodbvar, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

func dishupdate(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := dishes.Dishupdate(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

func dishdeletedisplay(httpwriter http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	// Get all selected records
	dishselected := req.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	type ControllerInfo struct {
		Name string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		DishItem   dishes.Dish
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishdelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Delete"

	items.DishItem = dishes.Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = dishes.Dish{}
	var dishname = items.DishItem.Name

	dishfind = dishes.Find(mongodbvar, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

func dishdelete(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := dishes.Dishdelete(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}
