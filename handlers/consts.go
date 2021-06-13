package handlers

const InternalServerError = "Internal Server error when fetching the data. Sorry for the inconveniance and try again later."

var REQUESTS_PER_HOUR = map[string]int{"free": 100, "basic": 1000, "lilleBoy": 100000, "GOD": -1}
var TIERS = []string{"free", "basic", "lilleBoy", "GOD"}
