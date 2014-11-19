package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/williambailey/consul-register/action"
)

var cmdApply = &Command{
	Usage: "apply [options] file.json",
	Short: "Apply a list of actions to the consul server.",
	Long: `
file.json contains an array of { "Action": "", "Config": {} } items that get applied in order.	
`,
	Run: runApply,
}

var (
	flagApply struct {
		server string
		token  string
		dry    bool
	}
)

func init() {
	consulFlag(&cmdApply.Flag, &flagApply.server, &flagApply.token)
	cmdApply.Flag.BoolVar(&flagApply.dry, "dry", false, "Perform a dry run.")
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
	if flagApply.dry {
		fmt.Println("!! Dry run.")
	}
	t := len(actions)
	f := fmt.Sprintf("%%%dd of %d - %%s\n", len(strconv.Itoa(t)), t)
	for i, a := range actions {
		fmt.Printf(f, i+1, a)
		if flagApply.dry {
			//...
		} else {
			err = a.Action(ctx)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}
