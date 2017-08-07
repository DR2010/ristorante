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

	db, err = sql.Open("mysql", "root:oculos18@/gufcdraws")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"

	fmt.Println("Running...")

	http.HandleFunc("/loginpage", loginPage)
	http.HandleFunc("/dishlist", dishlist)
	http.HandleFunc("/orderlist", orderlist)
	http.HandleFunc("/testtemplate", testtemplate)
	http.HandleFunc("/printparm", printparm)
	http.HandleFunc("/", root) // setting router rule

	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	http.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))

	err := http.ListenAndServe(":1515", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func root(httpwriter http.ResponseWriter, r *http.Request) {

	// create new template
	var listtemplate = `
		{{define "listtemplate"}}
	
		<div style="width:800px;">
			<div style="width:300px; float:left;">
				<p> 
            This is the best restaurante in the city. We are still coming up with ideas :-)
				</p>
				<p/>
			</div>
			<div style="width:300px; float:right;">
	            In this right panel we can add stuff also!
			</div>
		</div>

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

	var dishlist = []dishes.Dish{
		dishes.Dish{Type: "Main", Name: "Feijoada", Price: "200.00"},
		dishes.Dish{Type: "Main", Name: "Batatada", Price: "130.00"},
		dishes.Dish{Type: "Desert", Name: "Pudim de leite", Price: "50.00"},
		dishes.Dish{Type: "Entree", Name: "Almonds", Price: "25.00"},
	}

	items := DisplayTemplate{}
	items.Info.Name = "Dish List"

	// Set colum names
	items.FieldNames = make([]string, 3)
	items.FieldNames[0] = "ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Price"

	// Set rows to be displayed
	items.Rows = make([]Row, len(dishlist))

	for i := 0; i < len(dishlist); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, 3)
		items.Rows[i].Description[0] = dishlist[i].Type
		items.Rows[i].Description[1] = dishlist[i].Name
		items.Rows[i].Description[2] = dishlist[i].Price
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

func orderlist2(httpwriter http.ResponseWriter, req *http.Request) {

	var loadinplace = `
		{{define "loadinplace"}}
		<table style="width:100%">
		<tr>
			<th>Hard Coded Field 1</th>
			<th>{{ .Cfield2 }} </th>
			<th>{{ .Cfield3 }} </th>
			<th>{{ .Cfield4 }} </th>
		</tr>
		<tr>
			<td>Order List</td>
			<td>My Order List 1</td>
			<td>My Order more </td>
			<td>50.00</td>
		</tr>
		<tr>
			<td>Order Client 1</td>
			<td>My Order 1</td>
			<td>My Order more </td>
			<td>50.00</td>
		</tr>
		</table>
		{{end}}
		`

	// create new template
	t, _ := template.ParseFiles("indextemplate.tmpl")
	t, _ = t.Parse(loadinplace)

	err1 := t.Execute(httpwriter, map[string]string{
		"Cfield1": "Order #",
		"Cfield2": "Order Name",
		"Cfield3": "Order Description",
		"Cfield4": "Order Cost",
	})

	// err1 := t.ExecuteTemplate(httpwriter, "loadinplace", map[string]string{
	// 	"Cfield1": "Order #",
	// 	"Cfield2": "Order Name",
	// 	"Cfield3": "Order Description",
	// 	"Cfield4": "Order Cost",
	// })

	if err1 != nil {
		panic(err1)
	}
}

func dishadd(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Name = req.FormValue("dishname")
	dishtoadd.Price = req.FormValue("dishprice")

	ret := dishes.Dishadd(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		http.ServeFile(httpwriter, req, "success.html")
		return
	}
}

func testtemplate(httpwriter http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("test.html", "t1.tmpl", "t2.tmpl")
	if err != nil {
		fmt.Println(err)
	}
	items := struct {
		Name string
		City string
	}{
		Name: "MyName",
		City: "MyCity",
	}
	t.Execute(httpwriter, items)
}
