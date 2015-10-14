package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "Daisuke Fujita (@dtan4)"
	app.Email = "dtanshi45@gmail.com"
	app.Usage = "Simple CLI for manipulating AWS EC2"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
