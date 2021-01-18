package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gordonbondon/maxminddb-cidrs/pkg/cidrs"
)

type CountriesValue []cidrs.Country

func (c *CountriesValue) Set(s string) error {
	params := strings.Split(s, ":")
	code := params[0]

	country := cidrs.Country{ISOCode: code}

	if len(params) == 2 {
		country.Subdivisions = strings.Split(params[1], ",")
	} else if len(params) > 2 {
		return fmt.Errorf("wrong format, extra semicolon")
	}

	*c = append(*c, country)

	return nil
}

func (c *CountriesValue) String() string {
	return ""
}

func main() {
	var dbPathFlag string
	var countriesFlag CountriesValue

	flag.StringVar(&dbPathFlag, "dbpath", "GeoIP2-City.mmdb", dbPathFlagHelp)
	flag.Var(&countriesFlag, "country", countriesFlagHelp)
	flag.Parse()

	results, err := cidrs.List(&cidrs.ListOptions{
		DBPath:    dbPathFlag,
		Countries: countriesFlag,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, ip := range results {
		fmt.Println(ip)
	}
}

const (
	dbPathFlagHelp    = `Path to GeoIP2 mmdb file, requires detailed GeoIP2-City to use subdivisions`
	countriesFlagHelp = `ISO country code, with additional comma separated list of subdivisions after semicolon.
Ex:
	-country CH:GE,ZH"

Can be passed multiple times`
)
