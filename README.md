# ddcost

[![CI](https://github.com/winebarrel/ddcost/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/ddcost/actions/workflows/ci.yml)

A tool that shows a breakdown of Datadog costs in a table.

## Usage

```
Usage: ddcost --api-key=STRING --app-key=STRING

Flags:
  -h, --help                  Show context-sensitive help.
      --api-key=STRING        Datadog API key ($DD_API_KEY).
      --app-key=STRING        Datadog APP key ($DD_APP_KEY).
  -v, --view="summary"        Cost breakdown view (summary, sub-org).
  -o, --output="table"        Formatting style for output (table, tsv, json, csv).
  -s, --start-month=STRING    Cost beginning this month.
  -e, --end-month=STRING      Cost ending this month.
      --estimate              Get estimated cost.
      --version
```

```
$ export DD_API_KEY=...
$ export DD_APP_KEY=...
$ ddcost -v sub-org -s 2022-12
       ORG       |       PRODUCT       | CHARGE TYPE | 2022-12 | 2023-01 | 2023-02 | 2023-03 | 2023-04
-----------------+---------------------+-------------+---------+---------+---------+---------+----------
  organization1  | fargate_container   | committed   |       1 |       1 |       1 |       1 |       1
                 |                     | on_demand   |       2 |       2 |       2 |       2 |       2
                 |                     | total       |       3 |       3 |       3 |       3 |       3
                 | logs_indexed_15day  | committed   |       0 |       0 |       0 |       0 |       0
                 |                     | on_demand   |    0.50 |    0.50 |    0.50 |    0.50 |    0.50
                 |                     | total       |    0.50 |    0.50 |    0.50 |    0.50 |    0.50
                 | total               |             |    3.50 |    3.50 |    3.50 |    3.50 |    3.50
                 |                     |             |         |         |         |         |
  organization2  | infra_host          | committed   |      10 |      10 |      10 |      10 |      10
                 |                     | on_demand   |      20 |      20 |      20 |      20 |      20
                 |                     | total       |      30 |      30 |      30 |      30 |      30
                 | logs_indexed_15day  | committed   |       1 |       1 |       1 |       1 |       1
                 |                     | on_demand   |    1.50 |    1.50 |    1.50 |    1.50 |    1.50
                 |                     | total       |    2.50 |    2.50 |    2.50 |    2.50 |    2.50
                 | total               |             |   32.50 |   32.50 |   32.50 |   32.50 |   32.50
```

## Installation

```
brew install winebarrel/ddcost/ddcost
```

## Related Links

- https://github.com/winebarrel/ddusage
