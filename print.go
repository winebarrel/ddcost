package ddcost

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/olekukonko/tablewriter"
	"github.com/winebarrel/ddcost/internal/util"
)

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
			byMonth[month] = *charge.Cost
		}
	}

	ms := []string{}

	for m := range monthSet {
		ms = append(ms, m)
	}

	sort.Strings(ms)
	return cbd, ms
}

func printTable(resp *datadogV2.CostByOrgResponse, out io.Writer) {
	cbd, months := breakdownCost(resp)

	table := tablewriter.NewWriter(out)
	table.SetBorder(false)

	header := []string{"org", "product", "charge_type"}
	header = append(header, months...)
	table.SetHeader(header)

	util.EachEntryWithSort(cbd, func(org string, costByProduct CostByProduct, iCostByProduct int) {
		if iCostByProduct != 0 {
			emptyLine := make([]string, len(months)+3)
			table.Append(emptyLine)
		}

		util.EachEntryWithSort(costByProduct, func(product string, costByChargeType CostByChargeType, iCostByProduct int) {
			if iCostByProduct != 0 {
				org = ""
			}

			util.EachEntryWithSort(costByChargeType, func(chargeType string, costByMonth CostByMonth, iCostByChargeType int) {
				row := []string{"", "", chargeType}

				if iCostByChargeType == 0 {
					row[0] = org
					row[1] = product
				}

				for _, m := range months {
					cost := costByMonth[m]
					row = append(row, fmt.Sprintf("%.2f", cost))
				}

				table.Append(row)
			})
		})
	})

	table.Render()
}

func printTSV(resp *datadogV2.CostByOrgResponse, out io.Writer) {
	cbd, months := breakdownCost(resp)

	header := []string{"org", "product", "charge_type"}
	header = append(header, months...)
	fmt.Fprintln(out, strings.Join(header, "\t"))

	util.EachEntryWithSort(cbd, func(org string, costByProduct CostByProduct, _ int) {
		util.EachEntryWithSort(costByProduct, func(product string, costByChargeType CostByChargeType, iCostByProduct int) {
			if iCostByProduct != 0 {
				org = ""
			}

			util.EachEntryWithSort(costByChargeType, func(chargeType string, costByMonth CostByMonth, iCostByChargeType int) {
				row := []string{"", "", chargeType}

				if iCostByChargeType == 0 {
					row[0] = org
					row[1] = product
				}

				for _, m := range months {
					cost := costByMonth[m]
					row = append(row, fmt.Sprintf("%.2f", cost))
				}

				fmt.Fprintln(out, strings.Join(row, "\t"))
			})
		})
	})
}

func printJSON(resp *datadogV2.CostByOrgResponse, out io.Writer) {
	cbd, _ := breakdownCost(resp)
	m, _ := json.MarshalIndent(cbd, "", "  ")
	fmt.Fprintln(out, string(m))
}
