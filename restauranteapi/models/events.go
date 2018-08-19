// Package models is a dish for packages
// -------------------------------------
// .../restauranteapi/models/events.go
// -------------------------------------
package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Event is an event like festa junina, feijoada or any event to sell stuff
type Event struct {
	SystemID    bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	ID          string        // ID of the event (FESTAJUNINA2018, FEIJOADAMAY2018)
	Name        string        // Name of the event (Festa Junina 2018, Feijoada ...)
	Description string        // Description of the event
	Location    string        // Location of the event
	Date        string        // Event Date
}
