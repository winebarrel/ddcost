package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/ddcost"
)

var version string

var cli struct {
	ddcost.Options
	Version kong.VersionFlag
}

func init() {
	log.SetFlags(0)
}

func main() {
	kong.Parse(
		&cli,
		kong.Vars{"version": version},
	)

	client, err := ddcost.NewClient(&cli.Options)

	if err != nil {
		log.Fatal(err)
	}

	err = client.PrintHistoricalCostByOrg(os.Stdout)

	if err != nil {
		log.Fatal(err)
	}
}
