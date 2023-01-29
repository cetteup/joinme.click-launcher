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
			validator := IPPortURLValidator{}

			// WHEN
			err := validator.Validate(givenURL)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPatternURLValidator(t *testing.T) {
	type test struct {
		name            string
		givenHost       string
		givenPattern    string
		wantErrContains string
	}

	tests := []test{
		{
			name:         "no error for valid game id",
			givenHost:    "1234567890",
			givenPattern: Frostbite3GameIdPattern,
		},
		{
			name:            "error for non-numeric game id",
			givenHost:       "not-a-game-id",
			givenPattern:    Frostbite3GameIdPattern,
			wantErrContains: "url hostname is not a valid game id",
		},
		{
			name:            "error for empty hostname",
			givenHost:       "",
			givenPattern:    Frostbite3GameIdPattern,
			wantErrContains: "url hostname is not a valid game id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			givenUrl := &url.URL{Host: tt.givenHost}
			validator := MakePatternURLValidator(tt.givenPattern)

			// WHEN
			err := validator.Validate(givenUrl)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSimpleCmdBuilder(t *testing.T) {
	type test struct {
		name                   string
		givenHost              string
		givenDefaultArgs       []string
		givenAppendDefaultArgs bool
		givenPrefixes          []string
		givenLaunchType        game_launcher.LaunchType
		expectedCmd            []string
	}

	tests := []test{
		{
			name:            "returns plain command if launch type is launch and join",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenLaunchType: game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:     []string{net.JoinHostPort("1.1.1.1", "16567")},
		},
		{
			name:            "returns prefixed command",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenLaunchType: game_launcher.LaunchTypeLaunchAndJoin,
			givenPrefixes:   []string{"+connect"},
			expectedCmd:     []string{"+connect", net.JoinHostPort("1.1.1.1", "16567")},
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
			name:            "returns nil slice if launch type is any but launch and join",
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
			builder := MakeSimpleCmdBuilder(tt.givenPrefixes...)

			// WHEN
			cmd, err := builder.GetArgs(mockRepository, u, config, tt.givenLaunchType)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCmd, cmd)
		})
	}
}

func TestOriginCmdBuilder(t *testing.T) {
	type test struct {
		name                   string
		givenHost              string
		givenDefaultArgs       []string
		givenAppendDefaultArgs bool
		givenOfferIDs          []string
		givenLaunchType        game_launcher.LaunchType
		expectedCmd            []string
	}

	tests := []test{
		{
			name:            "returns complete command if launch type is launch and join",
			givenHost:       "987654321",
			givenLaunchType: game_launcher.LaunchTypeLaunchAndJoin,
			givenOfferIDs:   []string{"123456"},
			expectedCmd:     []string{"origin2://game/launch?cmdParams=-gameMode%2520MP%2520-role%2520soldier%2520-asSpectator%2520false%2520-joinWithParty%2520false%2520-gameId%2520987654321&offerIds=123456"},
		},
		{
			name:            "returns command without game details if launch type is any but launch and join",
			givenHost:       "987654321",
			givenOfferIDs:   []string{"123456"},
			givenLaunchType: game_launcher.LaunchTypeLaunchOnly,
			expectedCmd:     []string{"origin2://game/launch?cmdParams=&offerIds=123456"},
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
			builder := MakeOriginCmdBuilder(tt.givenOfferIDs...)

			// WHEN
			cmd, err := builder.GetArgs(mockRepository, u, config, tt.givenLaunchType)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCmd, cmd)
		})
	}
}

func TestRefractorV1CmdBuilder(t *testing.T) {
	type test struct {
		name                   string
		givenHost              string
		givenQuery             string
		givenDefaultArgs       []string
		givenAppendDefaultArgs bool
		givenLaunchType        game_launcher.LaunchType
		expectedCmd            []string
	}

	tests := []test{
		{
			name:            "returns complete command if launch type is launch and join",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenLaunchType: game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:     []string{"+joinServer", "1.1.1.1", "+port", "16567"},
		},
		{
			name:            "adds mod argument if url contains mod param",
			givenHost:       net.JoinHostPort("1.1.1.1", "16567"),
			givenQuery:      "mod=xpack",
			givenLaunchType: game_launcher.LaunchTypeLaunchAndJoin,
			expectedCmd:     []string{"+joinServer", "1.1.1.1", "+port", "16567", "+game", "xpack"},
		},
		{
			name:            "returns nil slice if launch type is any but launch and join",
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
			u := &url.URL{Host: tt.givenHost, RawQuery: tt.givenQuery}
			config := game_launcher.Config{
				DefaultArgs:       tt.givenDefaultArgs,
				AppendDefaultArgs: tt.givenAppendDefaultArgs,
			}
			builder := RefractorV1CmdBuilder{}

			// WHEN
			cmd, err := builder.GetArgs(mockRepository, u, config, tt.givenLaunchType)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCmd, cmd)
		})
	}
}
