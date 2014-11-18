package main

import (
	"log"

	"github.com/williambailey/consul-register/action"
)

var cmdApply = &Command{
	Usage: "apply [options] file",
	Short: "Apply a list of actions to the consul server.",
	Long: `
blah
`,
	Run: runApply,
}

var (
	flagApply struct {
		server string
		token  string
	}
)

func init() {
	consulFlag(&cmdApply.Flag, &flagApply.server, &flagApply.token)
}

func runApply(cmd *Command, args []string) {
	var (
		err     error
		ctx     action.Ctx
		actions action.Actions
	)
	if len(args) != 1 {
		cmd.UsageExit(nil)
	}
	ctx.API, err = parseConsulFlag(flagApply.server, flagApply.token)
	if err != nil {
		cmd.UsageExit(err)
	}
	actions, err = loadJSONActions(args[0])
	if err != nil {
		cmd.UsageExit(err)
	}

	err = doApply(&ctx, actions)
	if err != nil {
		log.Fatalln(err)
	}
}

func doApply(ctx *action.Ctx, actions action.Actions) error {
	var err error
	for _, a := range actions {
		err = a.Action(ctx)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
