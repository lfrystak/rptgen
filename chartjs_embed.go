package rptgen

import _ "embed" // required for go:embed directives below

// chartJSSource holds the Chart.js 4.4.4 UMD bundle.
// To upgrade: download chart.umd.min.js from
// https://cdn.jsdelivr.net/npm/chart.js@<version>/dist/chart.umd.min.js
// into assets/ and update the version note above.
//
//go:embed assets/chart.umd.min.js
var chartJSSource string
