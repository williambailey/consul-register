package action

import (
	"fmt"

	api "github.com/armon/consul-api"
)

var (
	// DefaultFactories provides a global list of action factories
	DefaultFactories = make(Factories, 0)
)

// Factories provides a slice of ActionFactory items and a few utility funcs.
type Factories []Factory

// NewAction returns a new zeroed Actioner by id.
// The idea is that we can then uses these factories inside of any
// action loading function without having to duplicate lots of code.
func (fs Factories) NewAction(id string) (Actioner, error) {
	for _, f := range fs {
		a, err := f(id)
		if err == nil {
			return a, nil
		}
		if _, ok := err.(UnknownFactoryIDError); ok {
			continue
		}
		return nil, err
	}
	return nil, UnknownFactoryIDError(id)
}

// Factory provides a way to create an empty Actioner from an identifier.
type Factory func(id string) (Actioner, error)

// UnknownFactoryIDError is used by a ActionFactory when the identifier is not supported.
type UnknownFactoryIDError string

func (e UnknownFactoryIDError) Error() string {
	return fmt.Sprintf("Unknown action %q", string(e))
}

// Actions simply a list of Actioner instances that should be applied in order
type Actions []Actioner

// Actioner performs an action.
type Actioner interface {
	// Action performs the action using the provided context
	Action(c *Ctx) error
	// Validate check that the action is valid in its current state
	//Validate() error
	// String to give us a user friendly identifier for the actioner
	//String() string
}

//Ctx provides context information to the Actioner.
type Ctx struct {
	API *api.Client
}
