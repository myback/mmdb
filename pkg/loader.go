package mmdb

import (
	"sync"
)

func subLoader(wg *sync.WaitGroup, filename, keyName, valueName, valueNameBak string, out *ids) {
	defer wg.Done()

	if err := NewIDs(filename, keyName, valueName, valueNameBak, out); err != nil {
		panic(err)
	}
}

// LoadData load data from files
func LoadData(cityFilename, cityIDFilename, countryFilename, countryIDFilename string) (*ipdb, *ipdb) {
	wg := sync.WaitGroup{}

	cityIDs := ids{}
	wg.Add(1)
	go subLoader(&wg, cityIDFilename, "geoname_id", "city_name", "country_iso_code", &cityIDs)

	countryIDs := ids{}
	wg.Add(1)
	go subLoader(&wg, countryIDFilename, "geoname_id", "country_iso_code", "country_name", &countryIDs)

	wg.Wait()

	// fmt.Printf("%#v\n", countryIDs)
	// os.Exit(1)

	city := ipdb{}
	wg.Add(1)
	go func(filename string, in *ids, out *ipdb) {
		defer wg.Done()

		var err error
		err = NewDB(filename, "network", "geoname_id", in, out)
		if err != nil {
			panic(err)
		}
	}(cityFilename, &cityIDs, &city)

	country := ipdb{}
	wg.Add(1)
	go func(filename string, in *ids, out *ipdb) {
		defer wg.Done()

		var err error
		err = NewDB(filename, "network", "geoname_id", in, out)
		if err != nil {
			panic(err)
		}
	}(countryFilename, &countryIDs, &country)

	wg.Wait()

	cityIDs = ids{}
	countryIDs = ids{}

	return &country, &city
}
