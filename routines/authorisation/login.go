package authorisation

import (
	"fmt"
)

// Login Comments
//
func Login(user string, password string) string {
	fmt.Println("username:", user)
	fmt.Println("password:", password)

	// read user from database
	// compare the hash password
	// return ok or not to go
	// return a token

	return user
}

// ValidateToken comments
//
func ValidateToken(token string) string {

	var authchecked = "Yes"

	return authchecked
}
