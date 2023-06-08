package ddcost

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/araddon/dateparse"
)

type Client struct {
	options        *Options
	api            *datadogV2.UsageMeteringApi
	timeStartMonth time.Time
	timeEndMonth   time.Time
}

func NewClient(options *Options) (*Client, error) {
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewUsageMeteringApi(apiClient)

	client := &Client{
		options: options,
		api:     api,
	}

	if options.StartMonth != "" {
		t, err := dateparse.ParseAny(options.StartMonth)

		if err != nil {
			return nil, err
		}

		client.timeStartMonth = t
	} else {
		client.timeStartMonth = defaultStartMonth
	}

	if options.EndMonth != "" {
		t, err := dateparse.ParseAny(options.EndMonth)

		if err != nil {
			return nil, err
		}

		client.timeEndMonth = t
	} else {
		client.timeEndMonth = defaultEndMonth
	}

	return client, nil
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

func (client *Client) PrintHistoricalCostByOrg(out io.Writer) error {
	ctx := client.withAPIKey(context.Background())

	resp, _, err := client.api.GetHistoricalCostByOrg(
		ctx,
		client.timeStartMonth,
		*datadogV2.NewGetHistoricalCostByOrgOptionalParameters().
			WithEndMonth(client.timeEndMonth).
			WithView(client.options.View),
	)

	if err != nil {
		var dderr datadog.GenericOpenAPIError

		if errors.As(err, &dderr) {
			err = fmt.Errorf("%w: %s", err, dderr.ErrorBody)
		}

		return err
	}

	switch client.options.Output {
	case "table":
		printTable(&resp, out)
	case "tsv":
		printTSV(&resp, out)
	case "json":
		printJSON(&resp, out)
	}

	return nil
}
