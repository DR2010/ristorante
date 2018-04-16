// Package main is the main package
// -------------------------------------
// .../restauranteapi/securityhandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"net/http"
	"restauranteapi/security"
)

// Hsecuritylogin is
func Hsecuritylogin(httpwriter http.ResponseWriter, req *http.Request) {

	var userid = req.FormValue("userid")
	var password = req.FormValue("password")

	// params := req.URL.Query()
	// cotacaotoadd.Currency = params.Get("Currency")
	// cotacaotoadd.Balance = params.Get("Balance")

	token, _ := security.ValidateUserCredentials(redisclient, userid, password)

	if token == "Error" {
		httpwriter.WriteHeader(http.StatusInternalServerError)
		httpwriter.Write([]byte("500 - Something bad happened!"))
	}

	// Get user roles
	// Store jwt as key on cache
	// Store user roles also
	//
	var usercredentials security.Credentials
	usercredentials.UserID = userid
	// usercredentials.Roles = []string

	json.NewEncoder(httpwriter).Encode(&token)

}

// HsecurityloginV2 is
func HsecurityloginV2(httpwriter http.ResponseWriter, req *http.Request) {

	var userid = req.FormValue("userid")
	var password = req.FormValue("password")

	// params := req.URL.Query()
	// cotacaotoadd.Currency = params.Get("Currency")
	// cotacaotoadd.Balance = params.Get("Balance")

	credentialwithtoken, _ := security.ValidateUserCredentialsV2(redisclient, userid, password)

	if credentialwithtoken.JWT == "Error" {
		httpwriter.WriteHeader(http.StatusInternalServerError)
		httpwriter.Write([]byte("500 - Something bad happened!"))
	}

	json.NewEncoder(httpwriter).Encode(&credentialwithtoken)
}

// Hsecuritysignup is
func Hsecuritysignup(httpwriter http.ResponseWriter, req *http.Request) {

	var userInsert security.Credentials

	userInsert.UserID = req.FormValue("userid")
	userInsert.Name = req.FormValue("preferredname")
	userInsert.Password = req.FormValue("password")
	userInsert.PasswordValidate = req.FormValue("passwordvalidate")
	userInsert.ApplicationID = req.FormValue("applicationid")
	userInsert.CentroID = req.FormValue("centroid")
	userInsert.MobilePhone = req.FormValue("mobilephone")

	userInsert.ClaimSet = make([]security.Claim, 3)
	userInsert.ClaimSet[0].Type = "USERTYPE"
	userInsert.ClaimSet[0].Value = "BASIC"
	userInsert.ClaimSet[1].Type = "USERID"
	userInsert.ClaimSet[1].Value = userInsert.UserID
	userInsert.ClaimSet[2].Type = "APPLICATIONID"
	userInsert.ClaimSet[2].Value = req.FormValue("applicationid")

	token := ""
	_, resfind := security.Find(redisclient, userInsert.UserID)
	if resfind == "200 OK" {
		token = "User already exists"
	}

	// Add user
	results := security.Useradd(redisclient, userInsert)
	if results.ErrorCode == "200 OK" {
		token = results.ReturnedValue
	}

	json.NewEncoder(httpwriter).Encode(&token)

}
