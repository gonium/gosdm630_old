package main

import (
	"flag"
	"fmt"
	"github.com/gonium/gosdm630"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var rtuDevice = flag.String("rtuDevice", "/dev/ttyUSB0", "Path to serial RTU device")
var verbose = flag.Bool("verbose", false, "Enables extensive logging")
var interval = flag.Int("interval", 5, "Seconds between querying the SDM630")

func init() {
	flag.Parse()
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello foo")
}

func MkLastValueHandler(hc *sdm630.MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last := hc.GetLast()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := last.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func main() {
	var rc = make(sdm630.ReadingChannel)

	qe := sdm630.NewQueryEngine(*rtuDevice, *interval, *verbose, rc)
	go qe.Produce()
	hc := sdm630.NewMeasurementCache(rc, *interval)
	go hc.ConsumeData()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/last", MkLastValueHandler(hc))
	log.Fatal(http.ListenAndServe(":8080", router))
}
