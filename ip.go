package main

import (
	"flag"
	"net/http"

	"github.com/myback/mmdb/mmdb"
)

func main() {
	var cityFilename, cityIDFilename, countryFilename, countryIDFilename string

	flag.StringVar(&cityFilename, "cityDB", "GeoLite2-City/GeoLite2-City-Blocks-IPv4.csv", "Path to file GeoLite2-City-Blocks-IPv4.csv")
	flag.StringVar(&cityIDFilename, "cityLoc", "GeoLite2-City/GeoLite2-City-Locations-IPv4.csv", "Path to file GeoLite2-City-Locations-en.csv")
	flag.StringVar(&countryFilename, "countryDB", "GeoLite2-Country/GeoLite2-Country-Blocks-IPv4.csv", "Path to file GeoLite2-Country-Blocks-IPv4.csv")
	flag.StringVar(&countryIDFilename, "countryLoc", "GeoLite2-Country/GeoLite2-Country-Locations-IPv4.csv", "Path to file GeoLite2-Country-Blocks-IPv4.csv", "Path to file GeoLite2-Country-Locations-en.csv")
	flag.Parse()

	countryDB, cityDB := mmdb.LoadData(cityFilename, cityIDFilename, countryFilename, countryIDFilename)

	mux := http.NewServeMux()
	mux.HandleFunc("/city", get(cityDB))
	mux.HandleFunc("/country", get(countryDB))

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}

func get(db mmdb.Getter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(db.Get()))
	})
}
