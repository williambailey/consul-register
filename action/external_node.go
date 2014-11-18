package action

import (
	"errors"
	"fmt"

	api "github.com/armon/consul-api"
)

func init() {
	DefaultFactories = append(
		DefaultFactories,
		func(id string) (Actioner, error) {
			switch id {
			case "ExternalNodeRegister":
				return &ExternalNodeRegister{}, nil
			case "ExternalNodeDeregister":
				return &ExternalNodeDeregister{}, nil
			}
			return nil, UnknownFactoryIDError(id)
		},
	)
}

// ExternalNodeService holds information about a service provided by an external node
type ExternalNodeService struct {
	ID      string
	Service string
	Tags    []string
	Port    int
}

// ExternalNodeRegister action
type ExternalNodeRegister struct {
	Node     string
	Address  string
	Services []*ExternalNodeService
}

// Action registers an external node
func (a *ExternalNodeRegister) Action(c *Ctx) error {
	_, err := c.API.Catalog().Register(
		&api.CatalogRegistration{
			Node:    a.Node,
			Address: a.Address,
		},
		nil,
	)
	if err != nil {
		return err
	}
	for _, s := range a.Services {
		_, err := c.API.Catalog().Register(
			&api.CatalogRegistration{
				Node:    a.Node,
				Address: a.Address,
				Service: &api.AgentService{
					ID:      s.ID,
					Service: s.Service,
					Tags:    s.Tags,
					Port:    s.Port,
				},
			},
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate that the action is valid in its current state.
func (a *ExternalNodeRegister) Validate() error {
	if a.Node == "" {
		return errors.New("Node must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *ExternalNodeRegister) String() string {
	// TODO Better string for node vs. node service(s).
	return fmt.Sprintf("External Node Register %q", a.Node)
}

// ExternalNodeDeregister action
type ExternalNodeDeregister struct {
	Node     string
	Services []string
}

// Action deregisters an external node
func (a *ExternalNodeDeregister) Action(c *Ctx) error {
	if len(a.Services) < 1 {
		_, err := c.API.Catalog().Deregister(
			&api.CatalogDeregistration{
				Node: a.Node,
			},
			nil,
		)
		if err != nil {
			return err
		}
	} else {
		for _, s := range a.Services {
			_, err := c.API.Catalog().Deregister(
				&api.CatalogDeregistration{
					Node:      a.Node,
					ServiceID: s,
				},
				nil,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Validate that the action is valid in its current state.
func (a *ExternalNodeDeregister) Validate() error {
	if a.Node == "" {
		return errors.New("Node must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *ExternalNodeDeregister) String() string {
	// TODO Better string for full node vs. node service(s).
	return fmt.Sprintf("External Node Deregister %q", a.Node)
}
