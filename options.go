package ddcost

import (
	"time"

	"github.com/araddon/dateparse"
)

var (
	defaultStartMonth         time.Time
	defaultEndMonth           time.Time
	defaultEstimateStartMonth time.Time
)

func init() {
	defaultEndMonth = time.Now()
	defaultStartMonth = defaultEndMonth.AddDate(0, -6, 0)
	defaultEstimateStartMonth = defaultEndMonth.AddDate(0, -2, 0)
}

type ClientOptions struct {
	APIKey string `env:"DD_API_KEY" required:"" help:"Datadog API key."`
	APPKey string `env:"DD_APP_KEY" required:"" help:"Datadog APP key."`
}

type PrintHistoricalCostByOrgOptions struct {
	View       string `short:"v" enum:"summary,sub-org" default:"summary" help:"Cost breakdown view (summary, sub-org)."`
	Output     string `short:"o" enum:"table,tsv,json,csv" default:"table" help:"Formatting style for output (table, tsv, json, csv)."`
	StartMonth string `short:"s" help:"Cost beginning this month."`
	EndMonth   string `short:"e" help:"Cost ending this month."`
	Estimate   bool   `default:"false" help:"Get estimated cost."`
}

func (options *PrintHistoricalCostByOrgOptions) calcPeriod() (time.Time, time.Time, error) {
	timeStartMonth := defaultStartMonth
	timeEndMonth := defaultEndMonth

	if options.StartMonth != "" {
		t, err := dateparse.ParseAny(options.StartMonth)

		if err != nil {
			return timeStartMonth, timeEndMonth, err
		}

		timeStartMonth = t
	} else if options.Estimate {
		timeStartMonth = defaultEstimateStartMonth
	}

	if options.EndMonth != "" {
		t, err := dateparse.ParseAny(options.EndMonth)

		if err != nil {
			return timeStartMonth, timeEndMonth, err
		}

		timeEndMonth = t
	}

	return timeStartMonth, timeEndMonth, nil
}
