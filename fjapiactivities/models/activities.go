// Package models is a dish for packages
// -------------------------------------
// .../restauranteapi/models/dishes.go
// -------------------------------------
package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Activity is in fact an Event that happens (Feijoada, Picanha... etc)
type Activity struct {
	SystemID        bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	ShortName       string        // Short Name of the Event/ Activity (Feijoada10Jan2019, Picanha18May2020)
	Status          string        // Is the activity ACTIVE/ INACTIVE
	LongDescription string        // Long Name (feijoada dia 10 jan 2019, picanha dia 18 May 2020)
	Type            string        // PICKUPFOOD, ???
	StartDate       string        // Start Date
	EndDate         string        // End Date
	Description     string        // Description of the plate
}
