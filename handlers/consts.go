package handlers

import "math"

const InternalServerError = "Internal Server error when fetching the data. Sorry for the inconveniance and try again later."

var REQUESTS_PER_HOUR = map[string]float64{"free": 100, "basic": 1000, "lilleBoy": 100000, "GOD": math.Inf(1)}
var TIERS = []string{"free", "basic", "lilleBoy", "GOD"}
