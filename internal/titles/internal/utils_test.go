//go:build unit

package internal

import (
	"math"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEncryptedProfileConLogin(t *testing.T) {
	type test struct {
		name              string
		prepareProfileCon func(profileCon map[string]string)
		wantErrContains   string
	}

	tests := []test{
		{
			name:              "successfully extracts encrypted login details",
			prepareProfileCon: func(profileCon map[string]string) {},
		},
		{
			name: "fails if nickname is missing",
			prepareProfileCon: func(profileCon map[string]string) {
				delete(profileCon, ProfileNickConKey)
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if nickname is empty",
			prepareProfileCon: func(profileCon map[string]string) {
				profileCon[ProfileNickConKey] = ""
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if password is missing",
			prepareProfileCon: func(profileCon map[string]string) {
				delete(profileCon, ProfilePasswordConKey)
			},
			wantErrContains: "encrypted password is missing/empty",
		},
		{
			name: "fails if password is empty",
			prepareProfileCon: func(profileCon map[string]string) {
				profileCon[ProfilePasswordConKey] = ""
			},
			wantErrContains: "encrypted password is missing/empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			profileCon := map[string]string{
				ProfileNickConKey:     "mister249",
				ProfilePasswordConKey: "some-encrypted-password",
			}
			tt.prepareProfileCon(profileCon)

			// WHEN
			nickname, encryptedPassword, err := GetEncryptedProfileConLogin(profileCon)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, profileCon[ProfileNickConKey], nickname)
				assert.Equal(t, profileCon[ProfilePasswordConKey], encryptedPassword)
			}
		})
	}
}

func TestParseConFile(t *testing.T) {
	type test struct {
		name           string
		givenContent   string
		expectedResult map[string]string
	}

	tests := []test{
		{
			name:         "successfully parses .con file content",
			givenContent: "GlobalSettings.setDefaultUser \"0010\"\r\nGlobalSettings.setNamePrefix \"=DOG=\"\r\n",
			expectedResult: map[string]string{
				"GlobalSettings.setDefaultUser": "0010",
				"GlobalSettings.setNamePrefix":  "=DOG=",
			},
		},
		{
			name:         "concatenates multiple lines with the same key",
			givenContent: "GeneralSettings.addServerHistory \"135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 934\r\nGeneralSettings.addServerHistory \"138.197.130.124\" 29900 \"Weekend Warriors Wake Island\" 78",
			expectedResult: map[string]string{
				"GeneralSettings.addServerHistory": "135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 934,138.197.130.124\" 29900 \"Weekend Warriors Wake Island\" 78",
			},
		},
		{
			name:         "ignores lines not containing two space-separated elements",
			givenContent: "GeneralSettings.setSortOrder 0\r\nGeneralSettings.setNumRoundsPlayed",
			expectedResult: map[string]string{
				"GeneralSettings.setSortOrder": "0",
			},
		},
		{
			name:           "returns empty result for no content",
			givenContent:   "",
			expectedResult: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := ParseConFile([]byte(tt.givenContent))
			assert.Equal(t, tt.expectedResult, parsed)
		})
	}
}

func TestBuildOriginURL(t *testing.T) {
	type test struct {
		name          string
		givenOfferIDs []string
		givenArgs     []string
		expectedURL   string
	}

	tests := []test{
		{
			name:          "successfully builds URL with offer id and argument",
			givenOfferIDs: []string{"123"},
			givenArgs:     []string{"+launch"},
			expectedURL:   "origin2://game/launch?cmdParams=%2Blaunch&offerIds=123",
		},
		{
			name:          "successfully builds URL with multiple offer ids and args",
			givenOfferIDs: []string{"123", "456"},
			givenArgs:     []string{"+launch", "+quiet"},
			expectedURL:   "origin2://game/launch?cmdParams=%2Blaunch%2520%2Bquiet&offerIds=123%2C456",
		},
		{
			name:          "successfully builds URL without offer ids or args",
			givenOfferIDs: []string{},
			givenArgs:     []string{},
			expectedURL:   "origin2://game/launch?cmdParams=&offerIds=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originURL := BuildOriginURL(tt.givenOfferIDs, tt.givenArgs)
			assert.Equal(t, tt.expectedURL, originURL)
		})
	}
}

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
			name:        "true for for valid private IPv4",
			givenInput:  "192.168.1.1",
			wantIsValid: true,
		},
		{
			name:        "false for broadcast IPv4",
			givenInput:  "255.255.255.255",
			wantIsValid: false,
		},
		{
			name:        "false for IPv6",
			givenInput:  "2606:4700:4700::1111",
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
