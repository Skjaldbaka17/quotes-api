package handlers

import (
	"log"
	"math"
	"os"

	"github.com/joho/godotenv"
)

const InternalServerError = "Internal Server error when fetching the data. Sorry for the inconveniance and try again later."

var REQUESTS_PER_HOUR = map[string]float64{"free": 100, "basic": 1000, "lilleBoy": 100000, "GOD": math.Inf(1)}
var TIERS = []string{"free", "basic", "lilleBoy", "GOD"}

const DATABASE_URL = "DATABASE_URL"
const AUTHORS_TABLE = "AUTHORS_TABLE"
const QUOTES_TABLE = "QUOTES_TABLE"
const TOPICS_TABLE = "TOPICS_TABLE"
const USERS_TABLE = "USERS_TABLE"
const REQUEST_HISTORY_TABLE = "REQUEST_HISTORY_TABLE"
const AUTHOR_OF_THE_DAY_TABLE = "AUTHOR_OF_THE_DAY_TABLE"
const ICELANDIC_AUTHOR_OF_THE_DAY_TABLE = "ICELANDIC_AUTHOR_OF_THE_DAY_TABLE"
const SEARCH_VIEW = "SEARCH_VIEW"
const AUTHOR_OF_THE_DAY_VIEW = "AUTHOR_OF_THE_DAY_VIEW"
const ICELANDIC_AUTHOR_OF_THE_DAY_VIEW = "ICELANDIC_AUTHOR_OF_THE_DAY_VIEW"
const QUOTE_OF_THE_DAY_VIEW = "QUOTE_OF_THE_DAY_VIEW"
const ICELANDIC_QUOTE_OF_THE_DAY_VIEW = "ICELANDIC_QUOTE_OF_THE_DAY_VIEW"
const TOPICS_VIEW = "TOPICS_VIEW"
const ICELANDIC_QUOTE_OF_THE_DAY_TABLE = "ICELANDIC_QUOTE_OF_THE_DAY_TABLE"
const QUOTE_OF_THE_DAY_TABLE = "QUOTE_OF_THE_DAY_TABLE"

func GetEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load("./.env")

	if err != nil {
		log.Printf("Error loading .env file")
		err = godotenv.Load("../.env")
		if err != nil {
			log.Printf("Error loading ../.env file")
		}
	}

	return os.Getenv(key)
}
