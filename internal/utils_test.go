//go:build unit

package internal

import (
	"math"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryHasMod(t *testing.T) {
	type test struct {
		name       string
		givenQuery url.Values
		wantHasMod bool
	}

	tests := []test{
		{
			name: "true if query contains mod key",
			givenQuery: map[string][]string{
				urlQueryKeyMod: {"xpack"},
			},
			wantHasMod: true,
		},
		{
			name: "false if query does not contain mod key",
			givenQuery: map[string][]string{
				"some-other-query-param": {"some-value"},
			},
			wantHasMod: false,
		},
		{
			name:       "false if query is empty map",
			givenQuery: map[string][]string{},
			wantHasMod: false,
		},
		{
			name:       "false if query is nil",
			givenQuery: nil,
			wantHasMod: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasMod := QueryHasMod(tt.givenQuery)
			assert.Equal(t, tt.wantHasMod, hasMod)
		})
	}
}

func TestGetModFromQuery(t *testing.T) {
	type test struct {
		name        string
		givenQuery  url.Values
		expectedMod string
	}

	tests := []test{
		{
			name: "successfully extracts mod from query",
			givenQuery: map[string][]string{
				urlQueryKeyMod: {"xpack"},
			},
			expectedMod: "xpack",
		},
		{
			name: "empty string if query does not contain mod key",
			givenQuery: map[string][]string{
				"some-other-query-param": {"some-value"},
			},
			expectedMod: "",
		},
		{
			name:        "empty string if query is empty map",
			givenQuery:  map[string][]string{},
			expectedMod: "",
		},
		{
			name:        "empty string if query is nil",
			givenQuery:  nil,
			expectedMod: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := GetModFromQuery(tt.givenQuery)
			assert.Equal(t, tt.expectedMod, mod)
		})
	}
}

func TestIsValidIPv4(t *testing.T) {
	type test struct {
		name        string
		givenInput  string
		wantIsValid bool
	}

	tests := []test{
		{
			name:        "true for valid public IPv4",
			givenInput:  "1.1.1.1",
			wantIsValid: true,
		},
		{
			name:        "true for valid private IPv4",
			givenInput:  "192.168.1.1",
			wantIsValid: true,
		},
		{
			name:        "true for valid loopback IPv4",
			givenInput:  "127.0.0.1",
			wantIsValid: true,
		},
		{
			name:        "false for broadcast IPv4",
			givenInput:  "255.255.255.255",
			wantIsValid: false,
		},
		{
			name:        "false for public IPv6",
			givenInput:  "2606:4700:4700::1111",
			wantIsValid: false,
		},
		{
			name:        "false for private IPv6",
			givenInput:  "fd77:5a47:30c1:37e4::aaaa",
			wantIsValid: false,
		},
		{
			name:        "false for loopback IPv6",
			givenInput:  "::1",
			wantIsValid: false,
		},
		{
			name:        "false for invalid IPv4",
			givenInput:  "9.9.9.",
			wantIsValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := IsValidIPv4(tt.givenInput)
			assert.Equal(t, tt.wantIsValid, isValid)
		})
	}
}

func TestIsValidPort(t *testing.T) {
	type test struct {
		name        string
		givenPort   string
		wantIsValid bool
	}

	tests := []test{
		{
			name:        "true for minimum port",
			givenPort:   strconv.Itoa(portMin),
			wantIsValid: true,
		},
		{
			name:        "true for maximum port",
			givenPort:   strconv.Itoa(portMax),
			wantIsValid: true,
		},
		{
			name:        "false for port below minimum",
			givenPort:   strconv.Itoa(portMin - 1),
			wantIsValid: false,
		},
		{
			name:        "false for port above maximum",
			givenPort:   strconv.Itoa(portMax + 1),
			wantIsValid: false,
		},
		{
			name:        "false for port above int32 max",
			givenPort:   strconv.Itoa(math.MaxInt32 + 1),
			wantIsValid: false,
		},
		{
			name:        "false for non-numeric input",
			givenPort:   "not-a-port",
			wantIsValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := IsValidPort(tt.givenPort)
			assert.Equal(t, tt.wantIsValid, isValid)
		})
	}
}
