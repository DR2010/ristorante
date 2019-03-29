// Package main is the main package
// -------------------------------------
// .../restauranteapi/securityhandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"fjapisecurity/helper"
	"fjapisecurity/security"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Hsecuritylogin is
func Hsecuritylogin(httpwriter http.ResponseWriter, req *http.Request) {

	var userid = req.FormValue("userid")
	var password = req.FormValue("password")

	// params := req.URL.Query()
	// cotacaotoadd.Currency = params.Get("Currency")
	// cotacaotoadd.Balance = params.Get("Balance")

	token, _ := security.ValidateUserCredentials(userid, password)

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

	credentialwithtoken, _ := security.ValidateUserCredentialsV2(userid, password)

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
	userInsert.Status = "ACTIVE"

	userInsert.ClaimSet = make([]security.Claim, 3)
	userInsert.ClaimSet[0].Type = "USERTYPE"
	userInsert.ClaimSet[0].Value = "BASIC"
	userInsert.ClaimSet[1].Type = "USERID"
	userInsert.ClaimSet[1].Value = userInsert.UserID
	userInsert.ClaimSet[2].Type = "APPLICATIONID"
	userInsert.ClaimSet[2].Value = req.FormValue("applicationid")

	token := ""
	_, resfind := security.Find(userInsert.UserID)
	if resfind == "200 OK" {
		token = "User already exists"

		// 14-11-2018 Add else statement below
		// There was no exit when the user was found

	} else {

		// Add user
		results := security.Useradd(sysid, redisclient, userInsert)
		if results.ErrorCode == "200 OK" {
			token = results.ReturnedValue
		}
	}
	json.NewEncoder(httpwriter).Encode(&token)

}

// HgenerateCodeForgotPassword generates an ID and saves to Redis
func HgenerateCodeForgotPassword(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	type dcResetPassword struct {
		Email          string // User ID/ Email
		ResetCode      string // Code sent via email
		NewPassword    string // New Password
		RetypePassword string // New Password
	}

	var objtoaction dcResetPassword
	err = json.Unmarshal(bodybyte, &objtoaction)

	credentials := security.Credentials{}
	credentials.UserID = strings.ToUpper(objtoaction.Email)
	credentials.ResetCode = objtoaction.ResetCode

	rand.Seed(time.Now().UTC().UnixNano())
	code := strconv.Itoa(rand.Intn(100000))
	log.Println("Code:" + code)

	credentials.ResetCode = code

	// Expiry in 10 minutes
	//
	expiration := time.Minute * time.Duration(10)

	keyuser := "ResetPassword" + credentials.UserID

	err = redisclient.Set(keyuser, code, 0).Err()
	if err != nil {
		log.Println(err)
	}

	redisclient.Expire(keyuser, expiration)

	helper.SendEmail(credentials.UserID, "Your temporary code is "+credentials.ResetCode)

}

// HchangePassword generates an ID and saves to Redis
func HchangePassword(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	type dcResetPassword struct {
		Email          string // User ID/ Email
		ResetCode      string // Code sent via email
		NewPassword    string // New Password
		RetypePassword string // New Password
	}

	var objtoaction dcResetPassword
	err = json.Unmarshal(bodybyte, &objtoaction)

	credentials := security.Credentials{}
	credentials.UserID = strings.ToUpper(objtoaction.Email)
	credentials.ResetCode = objtoaction.ResetCode
	credentials.Password = security.Hashstring(objtoaction.NewPassword)
	credentials.PasswordValidate = security.Hashstring(objtoaction.RetypePassword)

	// Check if code is still valid
	//
	keyuser := "ResetPassword" + credentials.UserID

	codeinredis, _ := redisclient.Get(keyuser).Result()

	finalres := "Password has been updated."
	if codeinredis != credentials.ResetCode {
		log.Println("Code has expired.")
		finalres = "Code has expired."
		json.NewEncoder(httpwriter).Encode(&finalres)
		return
	}

	// Update Password

	usercredentials, resfind := security.Find(credentials.UserID)
	if resfind == "200 OK" {
		// User exists
		// Update password

		usercredentials.Password = security.Hashstring(objtoaction.NewPassword)
		usercredentials.PasswordValidate = security.Hashstring(objtoaction.RetypePassword)

		resultado := security.Userupdate(usercredentials)
		if resultado.IsSuccessful == "Y" {
			log.Println("All good, in theory.")
			finalres = "Password has been updated."
		} else {
			log.Println("not good")
			finalres = "not good"
		}

	}

	json.NewEncoder(httpwriter).Encode(&finalres)
	return

}

// HgetUserDetails retrieve user details
func HgetUserDetails(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	type dcUserDetails struct {
		Email         string // User ID/ Email
		IsAdmin       string //
		ApplicationID string //
		Status        string //
		UserType      string //
	}

	var objtoaction dcUserDetails
	err = json.Unmarshal(bodybyte, &objtoaction)

	credentials := security.Credentials{}
	credentials.UserID = strings.ToUpper(objtoaction.Email)

	usercredentials, resfind := security.Find(credentials.UserID)
	if resfind == "200 OK" {
		// All good
	} else {
		usercredentials.ApplicationID = "User Not found"
	}

	log.Println("found: " + usercredentials.UserID)
	num := len(usercredentials.ClaimSet)
	log.Printf("len claimset %d", num)

	for i := 0; i < len(usercredentials.ClaimSet); i++ {
		log.Println("Type: " + usercredentials.ClaimSet[i].Type)
		log.Println("Value: " + usercredentials.ClaimSet[i].Value)

		if usercredentials.ClaimSet[i].Type == "USERTYPE" {
			if usercredentials.ClaimSet[i].Value == "ADMIN" {
				usercredentials.IsAdmin = "Y"
				break
			}
		}
	}
	log.Println("isadmin: " + usercredentials.IsAdmin)
	log.Println("status: " + usercredentials.Status)

	json.NewEncoder(httpwriter).Encode(&usercredentials)
	return

}

// HlistAllUsers retrieve user details
func HlistAllUsers(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	type dcUserDetails struct {
		Email         string // User ID/ Email
		IsAdmin       string //
		ApplicationID string //
		Status        string //
		UserType      string //
	}

	var objtoaction dcUserDetails
	err = json.Unmarshal(bodybyte, &objtoaction)

	credentials := security.Credentials{}
	credentials.UserID = strings.ToUpper(objtoaction.Email)

	usercredentials, resfind := security.UsersGetAll()
	if resfind == "200 OK" {
		// All good
	} else {
		// usercredentials.ApplicationID = "User Not found"
		return
	}

	json.NewEncoder(httpwriter).Encode(&usercredentials)
	return

}

// HupdateUserDetails generates an ID and saves to Redis
func HupdateUserDetails(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	type dcUserDetails struct {
		Email         string // User ID/ Email
		IsAdmin       string //
		ApplicationID string //
		Status        string //
		UserType      string //
	}

	var objtoaction dcUserDetails
	err = json.Unmarshal(bodybyte, &objtoaction)

	credentials := security.Credentials{}
	credentials.UserID = strings.ToUpper(objtoaction.Email)
	credentials.IsAdmin = objtoaction.IsAdmin
	credentials.ApplicationID = objtoaction.ApplicationID
	credentials.Status = objtoaction.Status

	finalres := "Details have been updated."

	// Update Password

	usercredentials, resfind := security.Find(credentials.UserID)
	if resfind == "200 OK" {
		// User exists
		// Update password

		usercredentials.UserID = strings.ToUpper(objtoaction.Email)
		usercredentials.IsAdmin = objtoaction.IsAdmin
		usercredentials.ApplicationID = objtoaction.ApplicationID
		usercredentials.Status = objtoaction.Status
		usercredentials.ClaimSet[0].Type = "USERTYPE"
		usercredentials.ClaimSet[0].Value = objtoaction.UserType

		resultado := security.Userupdate(usercredentials)
		if resultado.IsSuccessful == "Y" {
			log.Println("All good, in theory.")
			finalres = "Details have been updated."
		} else {
			log.Println("not good")
			finalres = "not good"
		}

	}

	json.NewEncoder(httpwriter).Encode(&finalres)
	return

}
