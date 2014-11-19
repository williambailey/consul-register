package main

import (
	"log"
	"os"

	api "github.com/armon/consul-api"
	"github.com/williambailey/consul-register/action"
)

var cmdExport = &Command{
	Usage: "export [options]",
	Short: "Export consul configuration.",
	Long: `
    Configuration is exported in JSON format and sent directly to STDOUT.
    `,
	Run: runExport,
}

var (
	flagExport struct {
		server       string
		token        string
		acl          bool
		externalNode bool
		kv           bool
	}
)

func init() {
	consulFlag(&cmdExport.Flag, &flagExport.server, &flagExport.token)
	cmdExport.Flag.BoolVar(&flagExport.acl, "acl", false, "Include ACL.")
	cmdExport.Flag.BoolVar(&flagExport.externalNode, "externalNode", false, "Include External Nodes.")
	cmdExport.Flag.BoolVar(&flagExport.kv, "kv", false, "Include KV.")
}

func runExport(cmd *Command, args []string) {
	var (
		err     error
		ctx     action.Ctx
		actions = make(action.Actions, 0)
	)
	if len(args) != 0 {
		cmd.UsageExit(nil)
	}
	ctx.API, err = parseConsulFlag(flagExport.server, flagExport.token)
	if err != nil {
		cmd.UsageExit(err)
	}
	if flagExport.acl {
		actions, err = exportACL(&ctx, actions)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if flagExport.externalNode {
		actions, err = exportExternalNode(&ctx, actions)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if flagExport.kv {
		actions, err = exportKV(&ctx, actions)
		if err != nil {
			log.Fatalln(err)
		}
	}
	out, err := saveJSONActions(actions)
	if err != nil {
		log.Fatalln(err)
	}
	out.WriteTo(os.Stdout)
}

func exportACL(ctx *action.Ctx, a action.Actions) (action.Actions, error) {
	var err error
	acls, _, err := ctx.API.ACL().List(nil)
	if err != nil {
		return nil, err
	}
	for _, acl := range acls {
		if acl.Type == api.ACLManagementType {
			continue
		}
		a = append(a, &action.ACLSet{
			Name:  acl.Name,
			Rules: acl.Rules,
		})
	}
	return a, nil
}

func exportExternalNode(ctx *action.Ctx, a action.Actions) (action.Actions, error) {
	var err error
	q := &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}
	nodes, _, err := ctx.API.Catalog().Nodes(q)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		exportNode := true
		node, _, err := ctx.API.Catalog().Node(n.Node, q)
		if err != nil {
			return nil, err
		}
		for _, s := range node.Services {
			if s.ID == "consul" {
				exportNode = false
			}
		}
		if !exportNode {
			continue
		}
		en := &action.ExternalNodeRegister{
			Node:     node.Node.Node,
			Address:  node.Node.Address,
			Services: make([]*action.ExternalNodeService, len(node.Services)),
		}
		o := -1
		for _, s := range node.Services {
			o++
			en.Services[o] = &action.ExternalNodeService{
				ID:      s.ID,
				Service: s.Service,
				Tags:    s.Tags,
				Port:    s.Port,
			}
		}
		a = append(a, en)
	}
	return a, nil
}

func exportKV(ctx *action.Ctx, a action.Actions) (action.Actions, error) {
	var err error
	kvs, _, err := ctx.API.KV().List("", nil)
	if err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		a = append(a, &action.KVSet{
			Key:   kv.Key,
			Flags: kv.Flags,
			Value: string(kv.Value),
		})
	}
	return a, nil
}
