package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"restaurante/mongodb/dishes"
	"restaurante/mongodb/drinks"

	"restaurante/mongodb/helper"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var mongodbvar helper.DatabaseX

var db *sql.DB
var err error

// --------------
// Main function
// --------------
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

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/dishadd", dishadd)
	http.HandleFunc("/dishlist", showavailabledishes)
	http.HandleFunc("/drinkadd", showavailabledrinks)
	http.HandleFunc("/placeordervariable", placeordervariable)
	http.HandleFunc("/placeorderfile", placeorderfile)
	http.HandleFunc("/handler", handler)
	http.HandleFunc("/index2", index2)
	http.HandleFunc("/", homePage)

	// http.Handle("/node_modules/", http.StripPrefix("/node_modules", http.FileServer(http.Dir("./node_modules"))))
	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	http.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))

	http.ListenAndServe(":8080", nil)
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
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	res.Write([]byte("Hello" + databaseUsername))

}

func showavailabledrinks(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "drinksadd.html")
		return
	}

	drinktoadd := drinks.Drinks{}

	drinktoadd.Type = req.FormValue("drinktype")
	drinktoadd.Name = req.FormValue("drinkname")
	drinktoadd.Price = req.FormValue("drinkprice")

	ret := drinks.Drinkadd(mongodbvar, drinktoadd)

	res.Write([]byte(ret))
	// http.Redirect(res, req, "/", 301)
}

func showavailabledishes(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.ServeFile(res, req, "dishadd.html")
		return
	}

	dishtoadd := dishes.Dish{}

	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Name = req.FormValue("dishname")
	dishtoadd.Price = req.FormValue("dishprice")

	ret := dishes.Dishadd(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		http.ServeFile(res, req, "success.html")
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("header.html", "footer.html")
	t.Execute(w, map[string]string{"Title": "My title", "Body": "Hi this is my body"})
}

func placeorderX(res http.ResponseWriter, req *http.Request) {

	http.ServeFile(res, req, "placeorder.html")

}

func placeordervariable(httpwriter http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("partials/view.html") //setp 1
	t.Execute(httpwriter, "Hello World!")             //step 2

}

func placeorderfile(httpwriter http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("partials/viewtemplate.html")        // Parse template file.
	t.ExecuteTemplate(httpwriter, "partials/viewtemplate.html", nil) // merge.
}

func index2(httpwriter http.ResponseWriter, r *http.Request) {
	// t, _ := template.ParseFiles("index2.html")
	t, _ := template.ParseFiles("index2.html", "dishadd.html")
	t.Execute(httpwriter, map[string]string{
		"Title":          "My title",
		"Body":           "Hi this is my body passed in value.",
		"DishPriceLabel": "Dish Price",
		"DishPriceValue": "100.00",
	})
}

func homePage(res http.ResponseWriter, req *http.Request) {

	http.ServeFile(res, req, "index.html")

	// t, _ := template.ParseFiles("index.html")
	// t.Execute(res, nil)

}

func dishadd(httpwriter http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		// http.ServeFile(res, req, "dishadd.html")

		t, _ := template.ParseFiles("index2.html", "dishadd.html")
		t.Execute(httpwriter, map[string]string{
			"Title":          "My title",
			"Body":           "Hi this is my body passed in value.",
			"DishPriceLabel": "Dish Price",
			"DishPriceValue": "100.00",
		})

		return
	}

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
