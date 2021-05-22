package handlers

import (
	"net/http"
)

func GetQuotesById(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Matur"))
}
