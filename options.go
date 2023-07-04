package ddcost

import (
	"time"
)

var (
	defaultStartMonth         time.Time
	defaultEndMonth           time.Time
	defaultEstimateStartMonth time.Time
)

func init() {
	now := time.Now()
	defaultEndMonth = now.AddDate(0, 0, -now.Day()+1)
	defaultStartMonth = defaultEndMonth.AddDate(0, -6, 0)
	defaultEstimateStartMonth = defaultEndMonth.AddDate(0, -2, 0)
}

type Options struct {
	APIKey     string `env:"DD_API_KEY" required:"" help:"Datadog API key."`
	APPKey     string `env:"DD_APP_KEY" required:"" help:"Datadog APP key."`
	View       string `short:"v" enum:"summary,sub-org" default:"summary" help:"Cost breakdown view (summary, sub-org)."`
	Output     string `short:"o" enum:"table,tsv,json,csv" default:"table" help:"Formatting style for output (table, tsv, json, csv)."`
	StartMonth string `short:"s" help:"Cost beginning this month."`
	EndMonth   string `short:"e" help:"Cost ending this month."`
	Estimate   bool   `default:"false" help:"Get estimated cost."`
}
