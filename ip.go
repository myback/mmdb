package main

import (
	"flag"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	mmdb "github.com/myback/mmdb/pkg"
)

const big = 0xFFFFFF

func main() {
	var cityFilename, cityIDFilename, countryFilename, countryIDFilename string

	flag.StringVar(&cityFilename, "cityDB", "GeoLite2-City/GeoLite2-City-Blocks-IPv4.csv", "Path to file GeoLite2-City-Blocks-IPv4.csv")
	flag.StringVar(&cityIDFilename, "cityLoc", "GeoLite2-City/GeoLite2-City-Locations-IPv4.csv", "Path to file GeoLite2-City-Locations-en.csv")
	flag.StringVar(&countryFilename, "countryDB", "GeoLite2-Country/GeoLite2-Country-Blocks-IPv4.csv", "Path to file GeoLite2-Country-Blocks-IPv4.csv")
	flag.StringVar(&countryIDFilename, "countryLoc", "GeoLite2-Country/GeoLite2-Country-Locations-IPv4.csv", "Path to file GeoLite2-Country-Locations-en.csv")
	flag.Parse()

	countryDB, cityDB := mmdb.LoadData(cityFilename, cityIDFilename, countryFilename, countryIDFilename)

	r := mux.NewRouter()
	r.HandleFunc("/city/{ip}", get(cityDB)).Methods("GET")
	r.HandleFunc("/country/{ip}", get(countryDB)).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 100 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func get(db mmdb.Getter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		if !IsIpv4(params["ip"]) {
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(db.Get(params["ip"])))
	})
}

func IsIpv4(s string) bool {
	var p [net.IPv4len]byte
	for i := 0; i < net.IPv4len; i++ {
		if len(s) == 0 {
			// Missing octets.
			return false
		}
		if i > 0 {
			if s[0] != '.' {
				return false
			}
			s = s[1:]
		}

		var n int
		var i int
		var ok bool
		for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
			n = n*10 + int(s[i]-'0')
			if n >= big {
				n = big
				ok = false
			}
		}

		if i == 0 {
			n = 0
			i = 0
			ok = false
		}

		ok = true

		if !ok || n > 0xFF {
			return false
		}

		s = s[i:]
		p[i] = byte(n)
	}
	if len(s) != 0 {
		return false
	}
	return true
}
