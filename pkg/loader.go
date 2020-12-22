package mmdb

import "sync"

// LoadData load data from files
func LoadData(cityFilename, cityIDFilename, countryFilename, countryIDFilename string) (*ipdb, *ipdb) {
	wg := sync.WaitGroup{}

	cityIDs := ids{}
	go func(filename string, out *ids) {
		wg.Add(1)
		defer wg.Done()

		var err error
		out, err = NewIDs(filename, "geoname_id", "city_name")
		if err != nil {
			panic(err)
		}
	}(cityIDFilename, &cityIDs)

	countryIDs := ids{}
	go func(filename string, out *ids) {
		wg.Add(1)
		defer wg.Done()

		var err error
		out, err = NewIDs(filename, "geoname_id", "country_name")
		if err != nil {
			panic(err)
		}
	}(countryIDFilename, &countryIDs)

	wg.Wait()

	city := ipdb{}
	go func(filename string, in *ids, out *ipdb) {
		wg.Add(1)
		defer wg.Done()

		var err error
		out, err = NewDB(filename, "network", "geoname_id", in)
		if err != nil {
			panic(err)
		}
	}(cityFilename, &cityIDs, &city)

	country := ipdb{}
	go func(filename string, in *ids, out *ipdb) {
		wg.Add(1)
		defer wg.Done()

		var err error
		out, err = NewDB(filename, "network", "registered_country_geoname_id", in)
		if err != nil {
			panic(err)
		}
	}(countryFilename, &countryIDs, &country)

	wg.Wait()

	cityIDs = ids{}
	countryIDs = ids{}

	return &country, &city
}
