// Package security is a security for packages
// -------------------------------------
// .../restauranteapi/security/security.go
// -------------------------------------
package security

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	helper "restauranteapi/helper"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Credentials is to be exported
type Credentials struct {
	SystemID         bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	UserID           string        //
	Name             string        //
	Password         string        //
	PasswordValidate string        //
	ApplicationID    string        //
	CentroID         string        //
	MobilePhone      string        //
	JWT              string        //
	ClaimSet         []Claim       //
	Status           string        // It is set to Active manually by Daniel 'Active' or Inactive.
	IsAdmin          string        // It is set to Active manually by Daniel 'Active' or Inactive.
}

// Claim is
type Claim struct {
	Type  string
	Value string
}

// Useradd is for export
func Useradd(redisclient *redis.Client, userInsert Credentials) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "security"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Insert(userInsert)

	if err != nil {
		log.Fatal(err)

		var resX helper.Resultado
		resX.ErrorCode = "405 Error creating user"
		resX.ErrorDescription = "Error adding user"
		resX.IsSuccessful = "N"
		return resX

	}

	var res helper.Resultado
	res.ErrorCode = "200 OK"
	res.ErrorDescription = "User added"
	res.IsSuccessful = "Y"

	return res
}

// Find is to find stuff
func Find(redisclient *redis.Client, userid string) (Credentials, string) {

	database := new(helper.DatabaseX)
	database.Collection = "security"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	dishnull := Credentials{}

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	result := []Credentials{}
	err1 := c.Find(bson.M{"userid": userid}).All(&result)
	if err1 != nil {
		log.Fatal(err1)
	}

	var numrecsel = len(result)

	if numrecsel <= 0 {
		return dishnull, "404 Not found"
	}

	return result[0], "200 OK"
}

// UsersGetAll is to find stuff
func UsersGetAll(redisclient *redis.Client, userid string) []Credentials {

	database := new(helper.DatabaseX)
	database.Collection = "security"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database.Database).C(database.Collection)

	var results []Credentials

	err1 := c.Find(nil).All(&results)
	if err1 != nil {
		log.Fatal(err1)
	}

	return results
}

// Userupdate is
func Userupdate(redisclient *redis.Client, userUpdate Credentials) helper.Resultado {

	database := new(helper.DatabaseX)
	database.Collection = "security"
	database.Database, _ = redisclient.Get("API.MongoDB.Database").Result()
	database.Location, _ = redisclient.Get("API.MongoDB.Location").Result()

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Update(bson.M{"name": userUpdate.UserID}, userUpdate)

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Something Happened"
	res.IsSuccessful = "Y"

	return res
}

// ValidateUserCredentials is to find stuff
func ValidateUserCredentials(redisclient *redis.Client, userid string, password string) (string, string) {

	// look for user
	var us, _ = Find(redisclient, userid)

	var passwordhashed = Hashstring(password)

	if passwordhashed != us.Password {
		return "Error", "404 Error"
	}

	var jwt = getjwtfortoday(userid)
	return jwt, "200 OK"
}

// ValidateUserCredentialsV2 is to find stuff
func ValidateUserCredentialsV2(redisclient *redis.Client, userid string, password string) (Credentials, string) {

	var usercredentials Credentials
	usercredentials.UserID = userid
	usercredentials.ApplicationID = "None"
	usercredentials.JWT = "Error"
	usercredentials.Status = "Error"

	// look for user
	var userdatabase, _ = Find(redisclient, userid)

	var passwordhashed = Hashstring(password)

	if passwordhashed != userdatabase.Password {
		usercredentials.Status = "404 Error invalid password"
		return usercredentials, "404 Error"
	}

	if userdatabase.ApplicationID == "Belnorth" {
		if userdatabase.Status != "Active" {
			// If I have not make the user active
			// initially only set to Belnorth
			usercredentials.Status = "404 Error Belnorth User not active"
			return usercredentials, "404 Error"
		}
	}

	// Get the JWT
	var jwt = getjwtfortoday(userid)

	// Assign the JWT to the return JSON object Credentials
	userdatabase.JWT = jwt

	userdatabase.IsAdmin = "No"
	// Check if user is admin
	for x := 0; x < len(userdatabase.ClaimSet); x++ {
		if userdatabase.ClaimSet[x].Type == "USERTYPE" {
			if userdatabase.ClaimSet[x].Value == "ADMIN" {
				// list all if user is admin
				userdatabase.IsAdmin = "Yes"
				break
			}
		}
	}

	return userdatabase, "200 OK"
}

func keyfortheday(day int) string {

	var key = "De tudo, ao meu amor serei atento antes" +
		"E com tal zelo, e sempre, e tanto" +
		"Que mesmo em face do maior encanto" +
		"Dele se encante mais meu pensamento" +
		"Quero vivê-lo em cada vão momento" +
		"E em seu louvor hei de espalhar meu canto" +
		"E rir meu riso e derramar meu pranto" +
		"Ao seu pesar ou seu contentamento" +
		"E assim quando mais tarde me procure" +
		"Quem sabe a morte, angústia de quem vive" +
		"Quem sabe a solidão, fim de quem ama" +
		"Eu possa lhe dizer do amor que tive" +
		"Que não seja imortal, posto que é chama" +
		"Mas que seja infinito enquanto dure"

	stringSlice := strings.Split(key, " ")
	var stringSliceFinal [100]string

	x := 0
	for i := 0; i < len(stringSlice); i++ {
		if len(stringSlice[i]) > 3 {
			stringSliceFinal[x] = stringSlice[i]
			x++
		}
	}

	return stringSliceFinal[day]
}

// getjwtfortoday
// this is just a reference key
// the roles, date and user will be stored at the server
func getjwtfortoday(user string) string {

	// Generate Key
	_, _, day := time.Now().Date()
	s := keyfortheday(day)
	s += user
	h := sha1.New()
	h.Write([]byte(s))

	sha1hash := hex.EncodeToString(h.Sum(nil))

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return sha1hash
}

// Hashstring is just to hash the string
func Hashstring(str string) string {

	s := str
	h := sha1.New()
	h.Write([]byte(s))

	sha1hash := hex.EncodeToString(h.Sum(nil))

	return sha1hash
}
