package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
)

type ConfigFile struct {
	Filename string
}

var args struct {
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
	rawArgs.Bind(&args)
	fmt.Println(args)
	config := ConfigFile{args.Filename}
	run_command(config)
}

func run_command(config ConfigFile) (err error) {
	fmt.Println(config)
	return nil
}
