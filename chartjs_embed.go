package rptgen

import _ "embed"

// ChartJSVersion is the version of the Chart.js UMD bundle embedded in this package.
// To upgrade: download the new chart.umd.min.js from
// https://cdn.jsdelivr.net/npm/chart.js@<version>/dist/chart.umd.min.js
// into assets/, then update this constant.
const ChartJSVersion = "4.4.4"

//go:embed assets/chart.umd.min.js
var chartJSSource string
