package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"text/template"
	"unicode"

	api "github.com/armon/consul-api"
	"github.com/williambailey/consul-register/action"
)

// A Command is an implementation of a consul-register command
// like consul-register foo or consul-register bar.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// Usage is the one-line usage message.
	// The first word in the line is taken to be the command name.
	Usage string

	// Short is the short description shown in the 'consul-register help' output.
	Short string

	// Long is the long message shown in the
	// 'consul-register help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

// Name extracts the name from the first word in c.Usage
func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// UsageDefaults get the flag default values as a string.
func (c *Command) UsageDefaults() string {
	b := bytes.NewBuffer([]byte(""))
	c.Flag.SetOutput(b)
	c.Flag.PrintDefaults()
	c.Flag.SetOutput(os.Stderr)
	return b.String()
}

// UsageExit displays usage for the command and then exist.
func (c *Command) UsageExit(msg interface{}) {
	fmt.Fprintf(os.Stderr, "Usage: consul-register %s\n\n", c.Usage)
	fmt.Fprintf(os.Stderr, "Run 'consul-register help %s' for help.\n", c.Name())
	if msg != nil {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", msg)
	}
	os.Exit(2)
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed
// by 'consul-register help'.
var commands = []*Command{
	cmdApply,
}

func main() {
	flag.Usage = usageExit
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("consul-register: ")
	args := flag.Args()
	if len(args) < 1 {
		usageExit()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.UsageExit(nil) }
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "consul-register: unknown command %q\n", args[0])
	fmt.Fprintf(os.Stderr, "Run 'consul-register help' for usage.\n")
	os.Exit(2)
}

var usageTemplate = `
consul-register is a tool for managing consul key value storage.
Usage:
  consul-register command [arguments]

The commands are:
{{range .}}
    {{.Name | printf "%-8s"}} {{.Short}}{{end}}

Use "consul-register help [command]" for more information about a command.
`

var helpTemplate = `
Usage: consul-register {{.Usage | trim}}
{{.UsageDefaults | trimRight}}

{{.Long | trim}}
`

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: consul-register help command\n\n")
		fmt.Fprintf(os.Stderr, "Too many arguments given.\n")
		os.Exit(2)
	}
	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}
}

func usageExit() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim":      strings.TrimSpace,
		"trimRight": func(s string) string { return strings.TrimRightFunc(s, unicode.IsSpace) },
	})
	template.Must(t.Parse(strings.TrimSpace(text) + "\n"))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func consulFlag(flag *flag.FlagSet, server, token *string) {
	flag.StringVar(server, "server", "http://127.0.0.1:8500", "Consul server address")
	flag.StringVar(token, "token", "", "Consul token")
}

func parseConsulFlag(consul, token string) (*api.Client, error) {
	// The api client wants scheme and address separately.
	var (
		address string
		scheme  string
		err     error
	)
	u, err := url.Parse(consul)
	if err != nil {
		return nil, fmt.Errorf("Invalid consul flag.\n\n%s", err)
	}
	if u.Scheme == "" {
		scheme = "http"
	} else {
		scheme = u.Scheme
	}
	u.Scheme = ""
	address = strings.TrimLeft(u.String(), "/")
	client, err := api.NewClient(
		&api.Config{
			Address: address,
			Scheme:  scheme,
			Token:   token,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to create consul api client.\n\n%s", err)
	}
	return client, nil
}

func loadJSONActions(filename string) (action.Actions, error) {
	type item struct {
		Action string
		Config json.RawMessage
	}
	var (
		actions action.Actions
		items   []item
		err     error
	)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %q.\n\n%s", file, err)
	}
	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		return nil, fmt.Errorf("Unable to load actions from JSON.\n\n%s", err)
	}
	for o, i := range items {
		a, err := action.DefaultFactories.NewAction(i.Action)
		if err != nil {
			return nil, fmt.Errorf("Unable to load action #%d, %s.\n\n%s", o+1, i.Action, err)
		}
		err = json.Unmarshal(i.Config, a)
		if err != nil {
			return nil, fmt.Errorf("Unable to load action #%d, %s.\n\n%s", o+1, i.Action, err)
		}
		actions = append(actions, a)
	}
	return actions, nil
}
