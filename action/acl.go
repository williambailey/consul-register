package action

import (
	api "github.com/armon/consul-api"
)

// ACLDelete action
type ACLDelete struct {
	Name string
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

// ACLSet action
type ACLSet struct {
	Name  string
	Rules string
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
