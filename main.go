package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
)

type Args struct {
	On       bool   `docopt:"on"`
	Off      bool   `docopt:"off"`
	Reset    bool   `docopt:"reset"`
	List     bool   `docopt:"list"`
	SetAlias bool   `docopt:"set-alias"`
	RmAlias  bool   `docopt:"rm-alias"`
	SetHost  bool   `docopt:"set-host"`
	Port     string `docopt:"<port>"`
	Name     string `docopt:"<name>"`
	Num      int    `docopt:"<num>"`
	Hostname string `docopt:"<hostname>"`
	Filename string `docopt:"--config"`
}

func main() {
	usage := `apc - Control APC network power strip

Usage:
  apc.py [options] (on [<port>] | off [<port>] | reset [<port>] | list)
  apc.py [options] set-alias <name> <num>
  apc.py [options] rm-alias <name>
  apc.py [options] set-host <hostname>
  apc.py --help

Commands:
  on                     Turn port on [defaults to last port if empty]
  off                    Turn port off [default to last port if empty]
  reset                  Reset port [default to last port if empty]
  list                   List all ports, their aliases, and their status
  set-alias              Set an alias for a port number
  rm-alias               Remove alias for a port
  set-host               Set host of APC device via IP address or hostname
  --help                 Print this usage screen

Options:
  --config <filename>    Point to custom config file [default: ~/.config/apc/config]`

	rawArgs, _ := docopt.ParseArgs(usage, nil, "1.0")
	var args Args
	rawArgs.Bind(&args)
	fmt.Println(args)
	config := NewConfigFile(args.Filename)
	runCommand(&args, config)
}

// --------------- Command Handlers ---------------

func onCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("on command")
	config.Read()
	config.Write()
	return nil
}

func offCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("off command")
	return nil
}

func resetCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("reset command")
	return nil
}

func listCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("list command")
	return nil
}

func setAliasCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("set alias command")
	return nil
}

func rmAliasCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("rm alias command")
	return nil
}

func setHostCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("set host command")
	return nil
}

func runCommand(args *Args, config *ConfigFile) (err error) {
	var bindings = []struct {
		key      bool
		callback func(*Args, *ConfigFile) error
	}{
		{args.On, onCommand},
		{args.Off, offCommand},
		{args.Reset, resetCommand},
		{args.List, listCommand},
		{args.SetAlias, setAliasCommand},
		{args.RmAlias, rmAliasCommand},
		{args.SetHost, setHostCommand},
	}
	for _, k := range bindings {
		if k.key {
			return k.callback(args, config)
		}
	}
	return errors.New("No valid command")
}
