package action

import (
	api "github.com/armon/consul-api"
)

// KVDelete action
type KVDelete struct {
	Key string
}

// Action performs the KV delete action
func (a *KVDelete) Action(c *Ctx) error {
	_, err := c.API.KV().Delete(a.Key, nil)
	return err
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

// KVSet action
type KVSet struct {
	Key   string
	Flags uint64
	Value []byte
}

// Action performs the KV set action
func (a *KVSet) Action(c *Ctx) error {
	p := &api.KVPair{
		Key:   a.Key,
		Flags: a.Flags,
		Value: a.Value,
	}
	_, err := c.API.KV().Put(p, nil)
	return err
}

// KVSetIfNotExist action
type KVSetIfNotExist struct {
	Key   string
	Flags uint64
	Value []byte
}

// Action performs the KV set if not exist action
func (a *KVSetIfNotExist) Action(c *Ctx) error {
	p := &api.KVPair{
		Key:   a.Key,
		Flags: a.Flags,
		Value: a.Value,
	}
	_, _, err := c.API.KV().CAS(p, nil)
	return err
}
