package ddcost

import (
	"time"
)

var (
	defaultStartMonth time.Time
	defaultEndMonth   time.Time
)

func init() {
	now := time.Now()
	defaultEndMonth = now.AddDate(0, 0, -now.Day()+1)
	defaultStartMonth = defaultEndMonth.AddDate(0, -6, 0)
}

type Options struct {
	APIKey     string    `env:"DD_API_KEY" required:"" help:"Datadog API key."`
	APPKey     string    `env:"DD_APP_KEY" required:"" help:"Datadog APP key."`
	View       string    `short:"v" enum:"summary,sub-org" default:"summary" help:"Cost breakdown view (summary, sub-org)."`
	Output     string    `short:"o" enum:"table,tsv,json" default:"table" help:"Formatting style for output (table, tsv, json)."`
	StartMonth time.Time `short:"s" help:"Cost beginning this month (default: half a year ago)."`
	EndMonth   time.Time `short:"e" help:"Cost ending this month (default: this month)."`
}
