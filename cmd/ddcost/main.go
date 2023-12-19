package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/ddcost"
)

var version string

type options struct {
	ddcost.ClientOptions
	ddcost.PrintHistoricalCostByOrgOptions
}

func init() {
	log.SetFlags(0)
}

func main() {
	var cli struct {
		options
		Version kong.VersionFlag
	}

	kong.Parse(
		&cli,
		kong.Vars{"version": version},
	)

	client := ddcost.NewClient(&cli.ClientOptions)
	err := client.PrintHistoricalCostByOrg(os.Stdout, &cli.PrintHistoricalCostByOrgOptions)

	if err != nil {
		log.Fatal(err)
	}
}
