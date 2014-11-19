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
			case "ACLDelete":
				return &ACLDelete{}, nil
			case "ACLSet":
				return &ACLSet{}, nil
			}
			return nil, UnknownFactoryIDError(id)
		},
	)
}

// ACLDelete action
type ACLDelete struct {
	Name string
}

// Type returns the type identifier for the actioner
func (a *ACLDelete) Type() string {
	return "ACLDelete"
}

// Action performs the ACL delete action
func (a *ACLDelete) Action(c *Ctx) error {
	q := &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}
	acls, _, err := c.API.ACL().List(q)
	if err != nil {
		return err
	}
	for _, acl := range acls {
		if acl.Name == a.Name {
			_, err = c.API.ACL().Destroy(acl.ID, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Validate that the action is valid in its current state.
func (a *ACLDelete) Validate() error {
	if a.Name == "" {
		return errors.New("Name must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *ACLDelete) String() string {
	return fmt.Sprintf("ACL Delete %q", a.Name)
}

// ACLSet action
type ACLSet struct {
	Name  string
	Rules string
}

// Type returns the type identifier for the actioner
func (a *ACLSet) Type() string {
	return "ACLSet"
}

// Action performs the ACL set
func (a *ACLSet) Action(c *Ctx) error {
	q := &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}
	acls, _, err := c.API.ACL().List(q)
	if err != nil {
		return err
	}
	updated := false
	for _, acl := range acls {
		if acl.Name == a.Name {
			_, err = c.API.ACL().Update(
				&api.ACLEntry{
					ID:    acl.ID,
					Name:  acl.Name,
					Type:  api.ACLClientType,
					Rules: a.Rules,
				},
				nil,
			)
			updated = true
			if err != nil {
				return err
			}
		}
	}
	if !updated {
		_, _, err = c.API.ACL().Create(
			&api.ACLEntry{
				Name:  a.Name,
				Type:  api.ACLClientType,
				Rules: a.Rules,
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
func (a *ACLSet) Validate() error {
	if a.Name == "" {
		return errors.New("Name must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *ACLSet) String() string {
	return fmt.Sprintf("ACL Set %q %q", a.Name, a.Rules)
}
