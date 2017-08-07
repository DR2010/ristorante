package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"mongodb/helper"
	"net/http"
	"routines/authorisation"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var mongodbvar helper.DatabaseX

var db *sql.DB
var err error

// Looks after the main routing
//
func main() {
	http.HandleFunc("/", anotherroot)         // setting router rule
	http.HandleFunc("/sayname", sayhelloName) // setting router rule
	http.HandleFunc("/parms", printparm)
	http.HandleFunc("/login", login)
	http.HandleFunc("/blog", blog)
	http.HandleFunc("/loginpage", loginPage)
	http.HandleFunc("/merge", testHandler)
	http.HandleFunc("/parse", parsefiles)
	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	http.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func anotherroot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
	return
}

func blog(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/blog.html")
	return
}

// Root selected
func root(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Write([]byte("path" + r.URL.Path))
	w.Write([]byte("scheme" + r.URL.Scheme))

	for k, v := range r.Form {
		w.Write([]byte("key:" + k))
		w.Write([]byte("val:" + strings.Join(v, "")))
	}

}

// How to access parms from the URL
func printparm(w http.ResponseWriter, r *http.Request) {

	// extracts parm values from URL
	//
	r.ParseForm()

	//Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	// fmt.Fprintf(w, "Hello astaxie!") // write data to response

	w.Write([]byte("path" + r.URL.Path))
	w.Write([]byte("scheme" + r.URL.Scheme))

	for k, v := range r.Form {
		w.Write([]byte("<key>:" + k))
		w.Write([]byte("<val>:" + strings.Join(v, "")))
	}
	w.Write([]byte("Hello astaxie!")) // write data to response

}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in

		// convert to string
		var u = strings.Join(r.Form["username"], "")
		var p = strings.Join(r.Form["password"], "")

		var ret = authorisation.Login(u, p)

		w.Write([]byte("Hello: " + ret))

	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	// fmt.Fprintf(w, "Hello astaxie!") // write data to response

	w.Write([]byte("path" + r.URL.Path))
	w.Write([]byte("scheme" + r.URL.Scheme))

	for k, v := range r.Form {
		w.Write([]byte("key:" + k))
		w.Write([]byte("val:" + strings.Join(v, "")))
	}
	w.Write([]byte("Hello astaxie!")) // write data to response

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

func testHandler(w http.ResponseWriter, r *http.Request) {
	//Parsing HTML
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
	t.Execute(w, items)
}

func parsefiles(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("layout.html")
	if err != nil {
		panic(err)
	}

	// err = t.ExecuteTemplate(os.Stdout, "filemerge.html", nil)
	err = t.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		panic(err)
	}

}
