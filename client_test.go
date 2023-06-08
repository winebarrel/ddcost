package ddcost_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/ddcost"
)

func TestSummaryTable(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "https://api.datadoghq.com/api/v2/usage/historical_cost", func(req *http.Request) (*http.Response, error) {
		assert.Equal("2023-01-01T00:00:00Z", req.URL.Query().Get("start_month"))
		assert.Equal("2023-03-01T00:00:00Z", req.URL.Query().Get("end_month"))
		assert.Equal("summary", req.URL.Query().Get("view"))

		response := httpmock.NewStringResponse(http.StatusOK, `
			{
				"data": [
					{
						"type": "cost_by_org",
						"id": "1",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-04-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "2",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-05-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					}
				]
			}
		`)

		return response, nil
	})

	options := &ddcost.Options{
		StartMonth: "2023/1",
		EndMonth:   "2023/3",
		View:       "summary",
		Output:     "table",
	}
	client, err := ddcost.NewClient(options)
	assert.NoError(err)

	var buf strings.Builder
	err = client.PrintHistoricalCostByOrg(&buf)
	assert.NoError(err)
	assert.Equal(strings.Join([]string{
		"   ORG   |      PRODUCT      | CHARGE TYPE | 2023-04 | 2023-05  ",
		"---------+-------------------+-------------+---------+----------",
		"  my-org | fargate_container | committed   |    2.00 |    2.00  ",
		"         |                   | on_demand   |    3.00 |    3.00  ",
		"         |                   | total       |    5.00 |    5.00  ",
		"         | infra_host        | committed   |    0.00 |    0.00  ",
		"         |                   | on_demand   |    1.00 |    1.00  ",
		"         |                   | total       |    1.00 |    1.00  ",
	}, "\n")+"\n", buf.String())
}

func TestSubOrgTable(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "https://api.datadoghq.com/api/v2/usage/historical_cost", func(req *http.Request) (*http.Response, error) {
		assert.Equal("2023-01-01T00:00:00Z", req.URL.Query().Get("start_month"))
		assert.Equal("2023-03-01T00:00:00Z", req.URL.Query().Get("end_month"))
		assert.Equal("sub-org", req.URL.Query().Get("view"))

		response := httpmock.NewStringResponse(http.StatusOK, `
			{
				"data": [
					{
						"type": "cost_by_org",
						"id": "1",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-04-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "2",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-05-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "3",
						"attributes": {
							"org_name": "my-org2",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-04-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "4",
						"attributes": {
							"org_name": "my-org2",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-05-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					}
				]
			}
		`)

		return response, nil
	})

	options := &ddcost.Options{
		StartMonth: "2023/1",
		EndMonth:   "2023/3",
		View:       "sub-org",
		Output:     "table",
	}
	client, err := ddcost.NewClient(options)
	assert.NoError(err)

	var buf strings.Builder
	err = client.PrintHistoricalCostByOrg(&buf)
	assert.NoError(err)
	assert.Equal(strings.Join([]string{
		"    ORG   |      PRODUCT      | CHARGE TYPE | 2023-04 | 2023-05  ",
		"----------+-------------------+-------------+---------+----------",
		"  my-org  | fargate_container | committed   |    2.00 |    2.00  ",
		"          |                   | on_demand   |    3.00 |    3.00  ",
		"          |                   | total       |    5.00 |    5.00  ",
		"          | infra_host        | committed   |    0.00 |    0.00  ",
		"          |                   | on_demand   |    1.00 |    1.00  ",
		"          |                   | total       |    1.00 |    1.00  ",
		"          |                   |             |         |          ",
		"  my-org2 | fargate_container | committed   |    2.00 |    2.00  ",
		"          |                   | on_demand   |    3.00 |    3.00  ",
		"          |                   | total       |    5.00 |    5.00  ",
		"          | infra_host        | committed   |    0.00 |    0.00  ",
		"          |                   | on_demand   |    1.00 |    1.00  ",
		"          |                   | total       |    1.00 |    1.00  ",
	}, "\n")+"\n", buf.String())
}

func TestSubOrgTSV(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "https://api.datadoghq.com/api/v2/usage/historical_cost", func(req *http.Request) (*http.Response, error) {
		assert.Equal("2023-01-01T00:00:00Z", req.URL.Query().Get("start_month"))
		assert.Equal("2023-03-01T00:00:00Z", req.URL.Query().Get("end_month"))
		assert.Equal("sub-org", req.URL.Query().Get("view"))

		response := httpmock.NewStringResponse(http.StatusOK, `
			{
				"data": [
					{
						"type": "cost_by_org",
						"id": "1",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-04-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "2",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-05-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "3",
						"attributes": {
							"org_name": "my-org2",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-04-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "4",
						"attributes": {
							"org_name": "my-org2",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-05-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					}
				]
			}
		`)

		return response, nil
	})

	options := &ddcost.Options{
		StartMonth: "2023/1",
		EndMonth:   "2023/3",
		View:       "sub-org",
		Output:     "tsv",
	}
	client, err := ddcost.NewClient(options)
	assert.NoError(err)

	var buf strings.Builder
	err = client.PrintHistoricalCostByOrg(&buf)
	assert.NoError(err)

	assert.Equal(`org	product	charge_type	2023-04	2023-05
my-org	fargate_container	committed	2.00	2.00
		on_demand	3.00	3.00
		total	5.00	5.00
	infra_host	committed	0.00	0.00
		on_demand	1.00	1.00
		total	1.00	1.00
my-org2	fargate_container	committed	2.00	2.00
		on_demand	3.00	3.00
		total	5.00	5.00
	infra_host	committed	0.00	0.00
		on_demand	1.00	1.00
		total	1.00	1.00
`, buf.String())
}

func TestSubOrgJSON(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "https://api.datadoghq.com/api/v2/usage/historical_cost", func(req *http.Request) (*http.Response, error) {
		assert.Equal("2023-01-01T00:00:00Z", req.URL.Query().Get("start_month"))
		assert.Equal("2023-03-01T00:00:00Z", req.URL.Query().Get("end_month"))
		assert.Equal("sub-org", req.URL.Query().Get("view"))

		response := httpmock.NewStringResponse(http.StatusOK, `
			{
				"data": [
					{
						"type": "cost_by_org",
						"id": "1",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-04-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "2",
						"attributes": {
							"org_name": "my-org",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-05-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "3",
						"attributes": {
							"org_name": "my-org2",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-04-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					},
					{
						"type": "cost_by_org",
						"id": "4",
						"attributes": {
							"org_name": "my-org2",
							"public_id": "1",
							"region": "us",
							"total_cost": 0.0,
							"date": "2023-05-01T00:00:00Z",
							"charges": [
								{
									"product_name": "infra_host",
									"charge_type": "committed",
									"cost": 0
								},
								{
									"product_name": "infra_host",
									"charge_type": "on_demand",
									"cost": 1.0
								},
								{
									"product_name": "infra_host",
									"charge_type": "total",
									"cost": 1.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "committed",
									"cost": 2.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "on_demand",
									"cost": 3.0
								},
								{
									"product_name": "fargate_container",
									"charge_type": "total",
									"cost": 5.0
								}
							]
						}
					}
				]
			}
		`)

		return response, nil
	})

	options := &ddcost.Options{
		StartMonth: "2023/1",
		EndMonth:   "2023/3",
		View:       "sub-org",
		Output:     "json",
	}
	client, err := ddcost.NewClient(options)
	assert.NoError(err)

	var buf strings.Builder
	err = client.PrintHistoricalCostByOrg(&buf)
	assert.NoError(err)

	assert.Equal(`{
  "my-org": {
    "fargate_container": {
      "committed": {
        "2023-04": 2,
        "2023-05": 2
      },
      "on_demand": {
        "2023-04": 3,
        "2023-05": 3
      },
      "total": {
        "2023-04": 5,
        "2023-05": 5
      }
    },
    "infra_host": {
      "committed": {
        "2023-04": 0,
        "2023-05": 0
      },
      "on_demand": {
        "2023-04": 1,
        "2023-05": 1
      },
      "total": {
        "2023-04": 1,
        "2023-05": 1
      }
    }
  },
  "my-org2": {
    "fargate_container": {
      "committed": {
        "2023-04": 2,
        "2023-05": 2
      },
      "on_demand": {
        "2023-04": 3,
        "2023-05": 3
      },
      "total": {
        "2023-04": 5,
        "2023-05": 5
      }
    },
    "infra_host": {
      "committed": {
        "2023-04": 0,
        "2023-05": 0
      },
      "on_demand": {
        "2023-04": 1,
        "2023-05": 1
      },
      "total": {
        "2023-04": 1,
        "2023-05": 1
      }
    }
  }
}
`, buf.String())
}
