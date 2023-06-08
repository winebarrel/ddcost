# ddcost

[![test](https://github.com/winebarrel/ddcost/actions/workflows/test.yml/badge.svg)](https://github.com/winebarrel/ddcost/actions/workflows/test.yml)

A tool that shows a breakdown of Datadog costs in a table.

## Usage

```
Usage: ddcost --api-key=STRING --app-key=STRING

Flags:
  -h, --help                  Show context-sensitive help.
      --api-key=STRING        Datadog API key ($DD_API_KEY).
      --app-key=STRING        Datadog APP key ($DD_APP_KEY).
  -v, --view="summary"        Cost breakdown view (summary, sub-org).
  -o, --output="table"        Formatting style for output (table, tsv, json).
  -s, --start-month=STRING    Cost beginning this month (default: half a year ago).
  -e, --end-month=STRING      Cost ending this month (default: this month).
      --version
```

```
$ export DD_API_KEY=...
$ export DD_APP_KEY=...
$ ddcost -v sub-org -s 2022-12
       ORG       |       PRODUCT       | CHARGE TYPE | 2022-12 | 2023-01 | 2023-02 | 2023-03 | 2023-04
-----------------+---------------------+-------------+---------+---------+---------+---------+----------
  organization1  | fargate_container   | committed   |    1.00 |    1.00 |    1.00 |    1.00 |    1.00
                 |                     | on_demand   |    2.00 |    2.00 |    2.00 |    2.00 |    2.00
                 |                     | total       |    3.00 |    3.00 |    3.00 |    3.00 |    3.00
                 | logs_indexed_15day  | committed   |    0.00 |    0.00 |    0.00 |    0.00 |    0.00
                 |                     | on_demand   |    0.50 |    0.50 |    0.50 |    0.50 |    0.50
                 |                     | total       |    0.50 |    0.50 |    0.50 |    0.50 |    0.50
                 |                     |             |         |         |         |         |
  organization2  | infra_host          | committed   |   10.00 |   10.00 |   10.00 |   10.00 |   10.00
                 |                     | on_demand   |   20.00 |   20.00 |   20.00 |   20.00 |   20.00
                 |                     | total       |   30.00 |   30.00 |   30.00 |   30.00 |   30.00
                 | logs_indexed_15day  | committed   |    1.00 |    1.00 |    1.00 |    1.00 |    1.00
                 |                     | on_demand   |    1.50 |    1.50 |    1.50 |    1.50 |    1.50
                 |                     | total       |    2.50 |    2.50 |    2.50 |    2.50 |    2.50
```

## Installation

```
brew install winebarrel/ddcost/ddcost
```
