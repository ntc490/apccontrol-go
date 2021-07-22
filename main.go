package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
)

// --------------- Program Usage ---------------

var usage = `apc - Control APC network power strip

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
	rawArgs, _ := docopt.ParseArgs(usage, nil, "1.0")
	var args Args
	rawArgs.Bind(&args)
	config := NewConfigFile(args.Filename)
	err := runCommand(&args, config)
	error := 0
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		error = -1
	}
	os.Exit(error)
}

// --------------- Command Handlers ---------------

func onCommand(args *Args, config *ConfigFile) (err error) {
	err = config.Read()
	if err != nil {
		return err
	}
	err = config.CheckBasicSettings()
	if err != nil {
		return err
	}
	apc := NewApcConnectionFromConfigFile(config)
	err = apc.On(args.Port)
	if err != nil {
		return err
	}
	// Update config.LastPort
	config.Write()
	return nil
}

func offCommand(args *Args, config *ConfigFile) (err error) {
	err = config.Read()
	if err != nil {
		return err
	}
	err = config.CheckBasicSettings()
	if err != nil {
		return err
	}
	apc := NewApcConnectionFromConfigFile(config)
	err = apc.Off(args.Port)
	if err != nil {
		return err
	}
	// Update config.LastPort
	config.Write()
	return nil
}

func resetCommand(args *Args, config *ConfigFile) (err error) {
	err = config.Read()
	if err != nil {
		return err
	}
	err = config.CheckBasicSettings()
	if err != nil {
		return err
	}
	apc := NewApcConnectionFromConfigFile(config)
	err = apc.Reset(args.Port)
	if err != nil {
		return err
	}
	// Update config.LastPort
	config.Write()
	return nil
}

func listCommand(args *Args, config *ConfigFile) (err error) {
	fmt.Println("list command")
	return nil
}

func setAliasCommand(args *Args, config *ConfigFile) (err error) {
	err = config.Read()
	if err != nil {
		return err
	}

	alias := args.Name
	num := args.Num
	os.Stderr.WriteString("Setting alias '" + alias + "'...\n")
	err = config.SetAlias(num, alias)
	if err == nil {
		config.Write()
	}
	return err
}

func rmAliasCommand(args *Args, config *ConfigFile) (err error) {
	alias := args.Name
	os.Stderr.WriteString("Removing alias '" + alias + "'...\n")
	err = config.Read()
	if err != nil {
		return err
	}
	err = config.RmAliasByName(alias)
	if err == nil {
		config.Write()
	}
	return err
}

func setHostCommand(args *Args, config *ConfigFile) (err error) {
	err = config.Read()
	if err != nil {
		return err
	}
	config.Hostname = args.Hostname
	config.Write()
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
