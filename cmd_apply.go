package main

import (
	"log"
	"time"

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
		err error
		ctx action.Ctx
	)
	ctx.API, err = parseConsulFlag(flagApply.server, flagApply.token)
	if err != nil {
		cmd.UsageExit(err)
	}

	err = doApply(&ctx, args)
	if err != nil {
		log.Fatalln(err)
	}
}

func doApply(ctx *action.Ctx, args []string) error {
	var err error
	a := action.KVSet{
		Key:   "test",
		Value: []byte(time.Now().Local().String()),
	}

	err = a.Action(ctx)
	if err != nil {
		panic(err)
	}

	return nil
}
