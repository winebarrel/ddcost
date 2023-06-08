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

type Cost float64

func (c Cost) String() string {
	return fmt.Sprintf("%.2f", c)
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

	emptyLine := make([]string, len(months)+3)

	for iCostByProduct, org := range util.MapSortKeys(cbd) {
		costByProduct := cbd[org]

		if iCostByProduct != 0 {
			table.Append(emptyLine)
		}

		for iCostByProduct, product := range util.MapSortKeys(costByProduct) {
			costByChargeType := costByProduct[product]

			if iCostByProduct != 0 {
				org = ""
			}

			for iCostByChargeType, chargeType := range util.MapSortKeys(costByChargeType) {
				costByMonth := costByChargeType[chargeType]
				row := []string{"", "", chargeType}

				if iCostByChargeType == 0 {
					row[0] = org
					row[1] = product
				}

				for _, m := range months {
					cost := costByMonth[m]
					row = append(row, cost.String())
				}

				table.Append(row)
			}
		}
	}

	table.Render()
}

func printTSV(resp *datadogV2.CostByOrgResponse, out io.Writer) {
	cbd, months := breakdownCost(resp)

	header := []string{"org", "product", "charge_type"}
	header = append(header, months...)
	fmt.Fprintln(out, strings.Join(header, "\t"))

	for _, org := range util.MapSortKeys(cbd) {
		costByProduct := cbd[org]

		for iCostByProduct, product := range util.MapSortKeys(costByProduct) {
			costByChargeType := costByProduct[product]

			if iCostByProduct != 0 {
				org = ""
			}

			for iCostByChargeType, chargeType := range util.MapSortKeys(costByChargeType) {
				costByMonth := costByChargeType[chargeType]
				row := []string{"", "", chargeType}

				if iCostByChargeType == 0 {
					row[0] = org
					row[1] = product
				}

				for _, m := range months {
					cost := costByMonth[m]
					row = append(row, cost.String())
				}

				fmt.Fprintln(out, strings.Join(row, "\t"))
			}
		}
	}
}

func printJSON(resp *datadogV2.CostByOrgResponse, out io.Writer) {
	cbd, _ := breakdownCost(resp)
	m, _ := json.MarshalIndent(cbd, "", "  ")
	fmt.Fprintln(out, string(m))
}
