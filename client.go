package ddcost

import (
	"context"
	"io"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type Client struct {
	options *Options
	api     *datadogV2.UsageMeteringApi
}

func NewClient(options *Options) *Client {
	if options.StartMonth.IsZero() {
		options.StartMonth = defaultStartMonth
	}

	if options.EndMonth.IsZero() {
		options.EndMonth = defaultEndMonth
	}

	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewUsageMeteringApi(apiClient)

	client := &Client{
		options: options,
		api:     api,
	}

	return client
}

func (client *Client) withAPIKey(ctx context.Context) context.Context {
	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: client.options.APIKey,
			},
			"appKeyAuth": {
				Key: client.options.APPKey,
			},
		},
	)

	return ctx
}

// charge_type/month/cost
type CostByMonth map[string]float64

// charge_type/month/cost
type CostByChargeType map[string]CostByMonth

// product_name/charge_type/month/cost
type CostByProduct map[string]CostByChargeType

// org_name/product_name/charge_type/month/cost
type CostBreakdown map[string]CostByProduct

func (client *Client) PrintHistoricalCostByOrg(out io.Writer) {
	ctx := client.withAPIKey(context.Background())

	resp, _, err := client.api.GetHistoricalCostByOrg(
		ctx,
		client.options.StartMonth,
		*datadogV2.NewGetHistoricalCostByOrgOptionalParameters().
			WithEndMonth(client.options.EndMonth).
			WithView(client.options.View),
	)

	if err != nil {
		panic(err)
	}

	switch client.options.Output {
	case "table":
		printTable(&resp, out)
	case "tsv":
		printTSV(&resp, out)
	case "json":
		printJSON(&resp, out)
	}
}
