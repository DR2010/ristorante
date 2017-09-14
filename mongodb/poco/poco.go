package poco

import (
	"gopkg.in/mgo.v2/bson"
)

// Order is to be exported
type Order struct {
	SystemID   bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	ID         string
	ClientName string
	Date       string
	Time       string
	Total      string
}
