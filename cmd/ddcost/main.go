package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/ddcost"
)

var version string

var cli struct {
	ddcost.Options
	Version kong.VersionFlag
}

func main() {
	kong.Parse(
		&cli,
		kong.Vars{"version": version},
	)

	client := ddcost.NewClient(&cli.Options)
	client.PrintHistoricalCostByOrg(os.Stdout)
}
