package ddcost

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/winebarrel/ddcost/internal/util"
)

type Cost float64

func (c Cost) String() string {
	return humanize.FtoaWithDigits(float64(c), 2)
}

func (c Cost) Float64() float64 {
	return float64(c)
}

// month/cost
type CostByMonth map[string]Cost

// charge_type/month/cost
type CostByChargeType map[string]CostByMonth

// product_name/charge_type/month/cost
type CostByProduct map[string]CostByChargeType

// org_name/product_name/charge_type/month/cost
type CostBreakdown map[string]CostByProduct

func breakdownCost(resp *datadogV2.CostByOrgResponse) (CostBreakdown, []string) {
	cbd := CostBreakdown{}
	monthSet := map[string]struct{}{}

	for _, data := range resp.Data {
		attrs := data.Attributes
		month := attrs.Date.Format("2006-01")
		monthSet[month] = struct{}{}
		byProduct := util.MapValueOrDefault(cbd, *attrs.OrgName, CostByProduct{})

		for _, charge := range attrs.Charges {
			byChargeType := util.MapValueOrDefault(byProduct, *charge.ProductName, CostByChargeType{})
			byMonth := util.MapValueOrDefault(byChargeType, *charge.ChargeType, CostByMonth{})
			byMonth[month] = Cost(*charge.Cost)
		}
	}

	return cbd, util.MapSortKeys(monthSet)
}

func printTable(resp *datadogV2.CostByOrgResponse, out io.Writer) {
	cbd, months := breakdownCost(resp)

	table := tablewriter.NewWriter(out)
	table.SetBorder(false)

	header := []string{"org", "product", "charge_type"}
	header = append(header, months...)
	table.SetHeader(header)

	printTable0(cbd, months, out, func(row []string) {
		table.Append(row)
	})

	table.Render()
}

func printTSV(resp *datadogV2.CostByOrgResponse, out io.Writer, sep string) {
	cbd, months := breakdownCost(resp)

	header := []string{"org", "product", "charge_type"}
	header = append(header, months...)
	fmt.Fprintln(out, strings.Join(header, sep))

	printTable0(cbd, months, out, func(row []string) {
		if strings.Join(row, "") != "" {
			fmt.Fprintln(out, strings.Join(row, sep))
		} else {
			fmt.Fprintln(out)
		}
	})
}

func printTable0(cbd CostBreakdown, months []string, out io.Writer, procRow func([]string)) {
	emptyLine := make([]string, len(months)+3)

	for idxOrg, org := range util.MapSortKeys(cbd) {
		costByProduct := cbd[org]
		orgTotal := map[string]float64{}

		if idxOrg != 0 {
			procRow(emptyLine)
		}

		for idxProduct, product := range util.MapSortKeys(costByProduct) {
			costByChargeType := costByProduct[product]

			if idxProduct != 0 {
				org = ""
			}

			for idxChargeType, chargeType := range util.MapSortKeys(costByChargeType) {
				costByMonth := costByChargeType[chargeType]
				row := []string{"", "", chargeType}

				if idxChargeType == 0 {
					row[0] = org
					row[1] = product
				}

				for _, m := range months {
					cost := costByMonth[m]
					row = append(row, cost.String())

					if chargeType == "total" {
						orgTotal[m] += cost.Float64()
					}
				}

				procRow(row)
			}
		}

		row := []string{"", "total", ""}

		for _, m := range months {
			cost := Cost(orgTotal[m])
			row = append(row, cost.String())
		}

		procRow(row)
	}
}

func printJSON(resp *datadogV2.CostByOrgResponse, out io.Writer) {
	cbd, _ := breakdownCost(resp)
	m, _ := json.MarshalIndent(cbd, "", "  ")
	fmt.Fprintln(out, string(m))
}
