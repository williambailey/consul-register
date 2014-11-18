package action

import (
	api "github.com/armon/consul-api"
)

//Ctx provides context information to the Actioner.
type Ctx struct {
	API *api.Client
}

// Actions simply a list of Actioner instances that should be applied in order
type Actions []Actioner

// Actioner performs an action.
type Actioner interface {
	Action(c *Ctx) error
}
