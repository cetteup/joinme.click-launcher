//go:build unit

package internal

import (
	"net"
	"net/url"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
)

func TestIPPortURLValidator(t *testing.T) {
	type test struct {
		name            string
		prepareURL      func(u *url.URL)
		wantErrContains string
	}

	tests := []test{
		{
			name:       "no error for url containing valid IPv4 and port",
			prepareURL: func(u *url.URL) {},
		},
		{
			name: "error for non IPv4 hostname",
			prepareURL: func(u *url.URL) {
				u.Host = "not-an-ipv4-address"
			},
			wantErrContains: "url hostname is not a valid IPv4 address",
		},
		{
			name: "error for empty port",
			prepareURL: func(u *url.URL) {
				u.Host = u.Hostname()
			},
			wantErrContains: "port is missing from url",
		},
		{
			name: "error for invalid port",
			prepareURL: func(u *url.URL) {
				u.Host = net.JoinHostPort(u.Hostname(), strconv.Itoa(65536))
			},
			wantErrContains: "url port is not a valid network port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			givenURL := &url.URL{
				Host: net.JoinHostPort("1.1.1.1", "16567"),
			}
			tt.prepareURL(givenURL)

			// WHEN
			err := IPPortURLValidator(givenURL)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFrostbite3GameIdURLValidator(t *testing.T) {
	type test struct {
		name            string
		givenGameID     string
		wantErrContains string
	}

	tests := []test{
		{
			name:        "no error for valid game id",
			givenGameID: "1234567890",
		},
		{
			name:            "error for non-numeric game id",
			givenGameID:     "not-a-game-id",
			wantErrContains: "url hostname is not a valid game id",
		},
		{
			name:            "error for empty hostname",
			givenGameID:     "",
			wantErrContains: "url hostname is not a valid game id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			givenUrl := &url.URL{Host: tt.givenGameID}

			// WHEN
			err := Frostbite3GameIdURLValidator(givenUrl)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPlusConnectCmdBuilder(t *testing.T) {
	type test struct {
		name                   string
		givenHost              string
		givenDefaultArgs       []string
		givenAppendDefaultArgs bool
		givenLaunchType        game_launcher.LaunchType
		expectedCmd            []string
	}

	tests := []test{
		{
			name:            "return connect command if launch type is launch and join",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenLaunchType: game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:     []string{"+connect", net.JoinHostPort("1.1.1.1", "16567")},
		},
		{
			name:             "appends to default arguments if any are given",
			givenHost:        net.JoinHostPort("1.1.1.1", "16567"),
			givenDefaultArgs: []string{"+launch"},
			givenLaunchType:  game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:      []string{"+launch", "+connect", net.JoinHostPort("1.1.1.1", "16567")},
		},
		{
			name:                   "prepends to default arguments if any are given and AppendDefaultArgs is true",
			givenHost:              net.JoinHostPort("1.1.1.1", "16567"),
			givenDefaultArgs:       []string{"+launch"},
			givenAppendDefaultArgs: true,
			givenLaunchType:        game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:            []string{"+connect", net.JoinHostPort("1.1.1.1", "16567"), "+launch"},
		},
		{
			name:            "return nil slice if launch type is any but launch and join",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenLaunchType: game_launcher.LaunchTypeLaunchOnly,
			expectedCmd:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			mockRepository := NewMockFileRepository(ctrl)
			u := &url.URL{Host: tt.givenHost}
			config := game_launcher.Config{
				DefaultArgs:       tt.givenDefaultArgs,
				AppendDefaultArgs: tt.givenAppendDefaultArgs,
			}

			// WHEN
			cmd, err := PlusConnectCmdBuilder(mockRepository, u, config, tt.givenLaunchType)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCmd, cmd)
		})
	}
}

func TestPlainCmdBuilder(t *testing.T) {
	type test struct {
		name                   string
		givenHost              string
		givenDefaultArgs       []string
		givenAppendDefaultArgs bool
		givenLaunchType        game_launcher.LaunchType
		expectedCmd            []string
	}

	tests := []test{
		{
			name:            "return plain command if launch type is launch and join",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenLaunchType: game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:     []string{net.JoinHostPort("1.1.1.1", "16567")},
		},
		{
			name:             "appends to default arguments if any are given",
			givenHost:        net.JoinHostPort("1.1.1.1", "16567"),
			givenDefaultArgs: []string{"+launch"},
			givenLaunchType:  game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:      []string{"+launch", net.JoinHostPort("1.1.1.1", "16567")},
		},
		{
			name:                   "prepends to default arguments if any are given and AppendDefaultArgs is true",
			givenHost:              net.JoinHostPort("1.1.1.1", "16567"),
			givenDefaultArgs:       []string{"+launch"},
			givenAppendDefaultArgs: true,
			givenLaunchType:        game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:            []string{net.JoinHostPort("1.1.1.1", "16567"), "+launch"},
		},
		{
			name:            "return nil slice if launch type is any but launch and join",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenLaunchType: game_launcher.LaunchTypeLaunchOnly,
			expectedCmd:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			mockRepository := NewMockFileRepository(ctrl)
			u := &url.URL{Host: tt.givenHost}
			config := game_launcher.Config{
				DefaultArgs:       tt.givenDefaultArgs,
				AppendDefaultArgs: tt.givenAppendDefaultArgs,
			}

			// WHEN
			cmd, err := PlainCmdBuilder(mockRepository, u, config, tt.givenLaunchType)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCmd, cmd)
		})
	}
}
