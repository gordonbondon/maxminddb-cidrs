# maxminddb-cidrs

Convenience command-line and library wrapper for [maxminddb-golang](https://github.com/oschwald/maxminddb-golang)
that helps retrieving network lists based on [ISO 3166-2](https://en.wikipedia.org/wiki/ISO_3166-2) country codes and subdivision codes.

## Usage

### CLI

```
Usage of maxminddb-cidrs:
  -country value
    	ISO country code, with additional comma separated list of subdivisions after semicolon.
    	Ex:
    		-country CH:GE,ZH"

    	Can be passed multiple times
  -dbpath string
    	Path to GeoIP2 mmdb file, requires detailed GeoIP2-City to use subdivisions (default "GeoIP2-City.mmdb")
  -ipv4
    	return only IPv4 networks
  -ipv6
    	return only IPv6 networks
```

### Library

```go
import (
		"fmt"

		"github.com/gordonbondon/maxminddb-cidrs/pkg/cidrs"
)

func main() {
	options := &cidrs.ListOptions{
		DBPath: "./GeoIP2-City.mmdb",
		IPv4: true,
		Countries: []cidrs.Country{
			{
				ISOCode:      "GB",
				Subdivisions: []string{"ENG"},
			},
		},
	}

	results, _ := cidrs.List(options)
	for _, ip := range results {
		fmt.Println(ip)
	}
}
```
