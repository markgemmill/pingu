package main

import (
	"github.com/alecthomas/kong"
	"pingu/pkg"
)

var console *pkg.Console

func main() {

	cli := &CLI{}
	ctx := kong.Parse(cli,
		kong.Name("pingu"),
		kong.Description("A url monitoring utility."),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Vars{
			"version": "0.1.0-dev.5",
		})

	err := ctx.Run(&Context{})
	pkg.ExitOnError(err, "")

}
