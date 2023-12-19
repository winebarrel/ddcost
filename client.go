package ddcost

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type Client struct {
	options *ClientOptions
	api     *datadogV2.UsageMeteringApi
}

func NewClient(options *ClientOptions) *Client {
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

func (client *Client) PrintHistoricalCostByOrg(out io.Writer, options *PrintHistoricalCostByOrgOptions) error {
	timeStartMonth, timeEndMonth, err := options.calcPeriod()

	if err != nil {
		return err
	}

	ctx := client.withAPIKey(context.Background())
	var resp datadogV2.CostByOrgResponse

	if options.Estimate {
		resp, _, err = client.api.GetEstimatedCostByOrg(
			ctx,
			*datadogV2.NewGetEstimatedCostByOrgOptionalParameters().
				WithStartMonth(timeStartMonth).
				WithEndMonth(timeEndMonth).
				WithView(options.View),
		)
	} else {
		resp, _, err = client.api.GetHistoricalCostByOrg(
			ctx,
			timeStartMonth,
			*datadogV2.NewGetHistoricalCostByOrgOptionalParameters().
				WithEndMonth(timeEndMonth).
				WithView(options.View),
		)
	}

	if err != nil {
		var dderr datadog.GenericOpenAPIError

		if errors.As(err, &dderr) {
			err = fmt.Errorf("%w: %s", err, dderr.ErrorBody)
		}

		return err
	}

	if len(resp.Data) == 0 {
		return errors.New("no data")
	}

	switch options.Output {
	case "table":
		printTable(&resp, out)
	case "tsv":
		printTSV(&resp, out, "\t")
	case "json":
		printJSON(&resp, out)
	case "csv":
		printTSV(&resp, out, ",")
	}

	return nil
}
