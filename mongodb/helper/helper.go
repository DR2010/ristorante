package helper

// DatabaseX is a struct
type DatabaseX struct {
	Location   string // location of the database localhost, something.com, etc
	Database   string // database name
	Collection string // collection name
	APIServer  string // collection name
}

// Resultado is a struct
type Resultado struct {
	ErrorCode        string // error code
	ErrorDescription string // description
	IsSuccessful     string // Y or N
	ReturnedValue    string // Any string
}

func add() {
}
