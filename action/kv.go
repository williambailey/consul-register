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
			case "KVDelete":
				return &KVDelete{}, nil
			case "KVDeleteTree":
				return &KVDeleteTree{}, nil
			case "KVSet":
				return &KVSet{}, nil
			case "KVSetIfNotExist":
				return &KVSetIfNotExist{}, nil
			}
			return nil, UnknownFactoryIDError(id)
		},
	)
}

// KVDelete action
type KVDelete struct {
	Key string
}

// Action performs the KV delete action
func (a *KVDelete) Action(c *Ctx) error {
	_, err := c.API.KV().Delete(a.Key, nil)
	return err
}

// Validate that the action is valid in its current state.
func (a *KVDelete) Validate() error {
	if a.Key == "" {
		return errors.New("Key must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *KVDelete) String() string {
	return fmt.Sprintf("KV Delete %q", a.Key)
}

// KVDeleteTree action
type KVDeleteTree struct {
	Prefix string
}

// Action performs the KV delete action
func (a *KVDeleteTree) Action(c *Ctx) error {
	_, err := c.API.KV().DeleteTree(a.Prefix, nil)
	return err
}

// Validate that the action is valid in its current state.
func (a *KVDeleteTree) Validate() error {
	if a.Prefix == "" {
		return errors.New("Prefix must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *KVDeleteTree) String() string {
	return fmt.Sprintf("KV Delete Tree %q", a.Prefix)
}

// KVSet action
type KVSet struct {
	Key   string
	Flags uint64
	Value string
}

// Action performs the KV set action
func (a *KVSet) Action(c *Ctx) error {
	p := &api.KVPair{
		Key:   a.Key,
		Flags: a.Flags,
		Value: []byte(a.Value),
	}
	_, err := c.API.KV().Put(p, nil)
	return err
}

// Validate that the action is valid in its current state.
func (a *KVSet) Validate() error {
	if a.Key == "" {
		return errors.New("Key must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *KVSet) String() string {
	return fmt.Sprintf("KV Set %q %d %q", a.Key, a.Flags, a.Value)
}

// KVSetIfNotExist action
type KVSetIfNotExist struct {
	Key   string
	Flags uint64
	Value string
}

// Action performs the KV set if not exist action
func (a *KVSetIfNotExist) Action(c *Ctx) error {
	p := &api.KVPair{
		Key:   a.Key,
		Flags: a.Flags,
		Value: []byte(a.Value),
	}
	_, _, err := c.API.KV().CAS(p, nil)
	return err
}

// Validate that the action is valid in its current state.
func (a *KVSetIfNotExist) Validate() error {
	if a.Key == "" {
		return errors.New("Key must not be empty.")
	}
	return nil
}

// String representation of the action.
func (a *KVSetIfNotExist) String() string {
	return fmt.Sprintf("KV Set If Not Exist %q %d %q", a.Key, a.Flags, a.Value)
}
