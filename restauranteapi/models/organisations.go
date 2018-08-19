// Package models is a organisation for packages
// -------------------------------------
// .../restauranteapi/models/organisation.go
// -------------------------------------
package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Organisation is an event like festa junina, feijoada or any event to sell stuff
type Organisation struct {
	SystemID      bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	ID            string        // ID of the organisation (CEHOTP, CESYSDNEY, )
	Name          string        // Name of the organisation (Centro Espirita The House of the Path, Centro Espirita de Sydney)
	Description   string        // Description of the event
	Date          string        // Data que a organisacao for registrada
	ContactPerson string        // Contact
	EmailAddress  string        // Contact email address
	ContactPhone  string        // Contact Phone Number
	Location      string        // City Name, State & Country
	Street        string        // Street Number and Name
}
