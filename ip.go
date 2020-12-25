package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	mmdb "github.com/myback/mmdb/pkg"
)

func main() {
	var addr, cityFilename, cityIDFilename, countryFilename, countryIDFilename string

	flag.StringVar(&addr, "addr", ":8080", "TCP address to listen to")
	flag.StringVar(&cityFilename, "cityDB", "GeoLite2-City/GeoLite2-City-Blocks-IPv4.csv", "Path to file GeoLite2-City-Blocks-IPv4.csv")
	flag.StringVar(&cityIDFilename, "cityLoc", "GeoLite2-City/GeoLite2-City-Locations-en.csv", "Path to file GeoLite2-City-Locations-en.csv")
	flag.StringVar(&countryFilename, "countryDB", "GeoLite2-Country/GeoLite2-Country-Blocks-IPv4.csv", "Path to file GeoLite2-Country-Blocks-IPv4.csv")
	flag.StringVar(&countryIDFilename, "countryLoc", "GeoLite2-Country/GeoLite2-Country-Locations-en.csv", "Path to file GeoLite2-Country-Locations-en.csv")
	flag.Parse()

	countryDB, cityDB := mmdb.LoadData(cityFilename, cityIDFilename, countryFilename, countryIDFilename)

	// p := profile.Start(profile.MemProfile, profile.ProfilePath("."))

	r := mux.NewRouter()
	r.HandleFunc("/city/{ip}", get(cityDB)).Methods("GET")
	r.HandleFunc("/country/{ip}", get(countryDB)).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 100 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	srv.Shutdown(ctx)

	// p.Stop()
}

func get(db mmdb.Getter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		w.WriteHeader(200)
		w.Write([]byte(db.Get(params["ip"])))
	})
}
