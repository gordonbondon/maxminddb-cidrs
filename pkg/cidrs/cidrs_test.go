package cidrs_test

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gordonbondon/maxminddb-cidrs/pkg/cidrs"
)

// TestNetrowksReader implements cidrs.NetworksReader
type TestNetrowksReader struct {
	networks []struct {
		Country      string
		Subdivisions []string
		IP           string
	}
	next struct {
		Country      string
		Subdivisions []string
		IP           string
	}
}

func (n *TestNetrowksReader) Next() bool {
	if len(n.networks) > 0 {
		n.next = n.networks[0]
		n.networks = n.networks[1:]

		return true
	}

	return false
}

func (n *TestNetrowksReader) Network(result interface{}) (*net.IPNet, error) {
	if n.next.Country != "" {
		network := struct {
			Country struct {
				IsoCode string
			}
			Subdivisions []struct {
				IsoCode string
			}
		}{
			Country: struct{ IsoCode string }{IsoCode: n.next.Country},
		}

		if len(n.next.Subdivisions) > 0 {
			network.Subdivisions = make([]struct{ IsoCode string }, len(n.next.Subdivisions))
			for i, s := range n.next.Subdivisions {
				network.Subdivisions[i] = struct{ IsoCode string }{IsoCode: s}
			}
		}

		// I'm so sorry
		value, _ := json.Marshal(network)
		json.Unmarshal(value, result)
	}

	_, ipnet, _ := net.ParseCIDR(n.next.IP)

	return ipnet, nil
}

func (n *TestNetrowksReader) Err() error {
	return nil
}

func TestCIDRs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		reader   *TestNetrowksReader
		options  cidrs.ListOptions
		expected []string
	}{
		{
			name: "matches ips for single country",
			reader: &TestNetrowksReader{
				networks: []struct {
					Country      string
					Subdivisions []string
					IP           string
				}{
					{
						Country: "GB",
						IP:      "192.0.2.1/32",
					},
					{
						Country: "NL",
						IP:      "192.0.2.2/32",
					},
					{
						Country: "GB",
						IP:      "192.0.2.3/32",
					},
				},
			},
			options: cidrs.ListOptions{
				Countries: []cidrs.Country{{ISOCode: "GB"}},
			},
			expected: []string{"192.0.2.1/32", "192.0.2.3/32"},
		},
		{
			name: "matches only subdivisions when",
			reader: &TestNetrowksReader{
				networks: []struct {
					Country      string
					Subdivisions []string
					IP           string
				}{
					{
						Country:      "GB",
						Subdivisions: []string{"ENG"},
						IP:           "192.0.2.1/32",
					},
					{
						Country: "GB",
						IP:      "192.0.2.2/32",
					},
					{
						Country:      "GB",
						Subdivisions: []string{"WLS"},
						IP:           "192.0.2.3/32",
					},
				},
			},
			options: cidrs.ListOptions{
				Countries: []cidrs.Country{{ISOCode: "GB", Subdivisions: []string{"ENG"}}},
			},
			expected: []string{"192.0.2.1/32"},
		},
		{
			name: "matches for multiple countries",
			reader: &TestNetrowksReader{
				networks: []struct {
					Country      string
					Subdivisions []string
					IP           string
				}{
					{
						Country:      "GB",
						Subdivisions: []string{"ENG"},
						IP:           "192.0.2.1/32",
					},
					{
						Country: "NL",
						IP:      "192.0.2.2/32",
					},
					{
						Country:      "GB",
						Subdivisions: []string{"WLS"},
						IP:           "192.0.2.3/32",
					},
				},
			},
			options: cidrs.ListOptions{
				Countries: []cidrs.Country{
					{ISOCode: "GB", Subdivisions: []string{"ENG"}},
					{ISOCode: "NL"},
				},
			},
			expected: []string{"192.0.2.1/32", "192.0.2.2/32"},
		},
		{
			name: "filters by ip type",
			reader: &TestNetrowksReader{
				networks: []struct {
					Country      string
					Subdivisions []string
					IP           string
				}{
					{
						Country: "GB",
						IP:      "192.0.2.1/32",
					},
					{
						Country: "GB",
						IP:      "2a02:f600::/29",
					},
				},
			},
			options: cidrs.ListOptions{
				IPv4:      true,
				Countries: []cidrs.Country{{ISOCode: "GB"}},
			},
			expected: []string{"192.0.2.1/32"},
		},
		{
			name: "filters by ip type ipv6",
			reader: &TestNetrowksReader{
				networks: []struct {
					Country      string
					Subdivisions []string
					IP           string
				}{
					{
						Country: "GB",
						IP:      "192.0.2.1/32",
					},
					{
						Country: "GB",
						IP:      "2a02:f600::/29",
					},
				},
			},
			options: cidrs.ListOptions{
				IPv6:      true,
				Countries: []cidrs.Country{{ISOCode: "GB"}},
			},
			expected: []string{"2a02:f600::/29"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.options.NetworksReader = tc.reader
			actual, _ := cidrs.List(&tc.options)

			assert.Equal(t, tc.expected, actual, "results should be equal")
		})
	}
}
