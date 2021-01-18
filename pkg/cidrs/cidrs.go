package cidrs

import (
	"net"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

// NetworksReader allows to replace *maxminddb.Networks
type NetworksReader interface {
	Next() bool
	Network(interface{}) (*net.IPNet, error)
	Err() error
}

// Country
type Country struct {
	ISOCode      string
	Subdivisions []string
}

// ListOptions
type ListOptions struct {
	// NetworksReader if not provided DBPath will be used to initialize mmdb reader
	NetworksReader NetworksReader

	// DBPath file path to mmdb file
	DBPath string

	// List of countries and subdivisions to select
	Countries []Country
}

// List will return list of CIDRs belonging to countries base on ListOptions provided
func List(options *ListOptions) ([]string, error) {
	if options.NetworksReader == nil {
		db, err := maxminddb.Open(options.DBPath)
		if err != nil {
			return nil, err
		}
		defer db.Close()

		options.NetworksReader = db.Networks(maxminddb.SkipAliasedNetworks)
	}

	match := make(map[string]map[string]int)

	for _, c := range options.Countries {
		match[c.ISOCode] = make(map[string]int)
		for _, s := range c.Subdivisions {
			match[c.ISOCode][s] = 0
		}
	}

	results := make([]string, 0, 10)

	record := struct {
		Country struct {
			IsoCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
		Subdivisions []struct {
			IsoCode string `maxminddb:"iso_code"`
		} `maxminddb:"subdivisions"`
	}{}

	networks := options.NetworksReader

	for networks.Next() {
		subnet, err := networks.Network(&record)
		if err != nil {
			return nil, err
		}

		if s, ok := match[record.Country.IsoCode]; ok {
			if len(s) != 0 {
				for _, i := range record.Subdivisions {
					if _, ok := s[i.IsoCode]; ok {
						results = append(results, subnet.String())
						continue
					}
				}
			} else {
				results = append(results, subnet.String())
			}
		}
	}
	if networks.Err() != nil {
		return nil, networks.Err()
	}

	return results, nil
}
