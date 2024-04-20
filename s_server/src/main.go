package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8000", New()))
}
