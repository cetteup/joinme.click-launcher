//go:build unit

package internal

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cetteup/joinme.click-launcher/internal/testhelpers"
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

func TestDeleteFileHookHandler(t *testing.T) {
	type test struct {
		name              string
		givenConfig       game_launcher.Config
		givenPathsBuilder func(config game_launcher.Config) ([]string, error)
		expect            func(fr *MockFileRepository)
		wantErrContains   string
	}

	tests := []test{
		{
			name: "deletes CoD running file if present in install path",
			givenConfig: game_launcher.Config{
				InstallPath:    "C:\\Program Files\\Call of Duty",
				ExecutableName: "CoDMP.exe",
			},
			givenPathsBuilder: CoDRunningFilePathsBuilder,
			expect: func(fr *MockFileRepository) {
				path := "C:\\Program Files\\Call of Duty\\__CoDMP"
				alternate := "AppData\\Local\\VirtualStore\\Program Files\\Call of Duty\\__CoDMP"
				fr.EXPECT().FileExists(gomock.Eq(path)).Return(true, nil)
				fr.EXPECT().RemoveAll(gomock.Eq(path)).Return(nil)
				fr.EXPECT().FileExists(testhelpers.StringContainsMatcher(alternate)).Return(false, nil)
			},
		},
		{
			name: "deletes CoD running file if present in virtual store",
			givenConfig: game_launcher.Config{
				InstallPath:    "C:\\Program Files (x86)\\Call of Duty 2",
				ExecutableName: "CoD2MP_s.exe",
			},
			givenPathsBuilder: CoDRunningFilePathsBuilder,
			expect: func(fr *MockFileRepository) {
				primary := "C:\\Program Files (x86)\\Call of Duty 2\\__CoD2MP_s"
				alternate := "AppData\\Local\\VirtualStore\\Program Files (x86)\\Call of Duty 2\\__CoD2MP_s"
				fr.EXPECT().FileExists(gomock.Eq(primary)).Return(false, nil)
				fr.EXPECT().FileExists(testhelpers.StringContainsMatcher(alternate)).Return(true, nil)
				fr.EXPECT().RemoveAll(testhelpers.StringContainsMatcher(alternate)).Return(nil)
			},
		},
		{
			name: "does nothing if running file does not exist",
			givenConfig: game_launcher.Config{
				InstallPath:    "C:\\Program Files\\Publisher\\Game",
				ExecutableName: "Game.exe",
			},
			givenPathsBuilder: func(config game_launcher.Config) ([]string, error) {
				return []string{filepath.Join(config.InstallPath, "Game.running")}, nil
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().FileExists(gomock.Eq("C:\\Program Files\\Publisher\\Game\\Game.running")).Return(false, nil)
			},
		},
		{
			name: "errors if paths builder fails",
			givenConfig: game_launcher.Config{
				InstallPath:    "C:\\Program Files\\Publisher\\Game",
				ExecutableName: "Game.exe",
			},
			givenPathsBuilder: func(config game_launcher.Config) ([]string, error) {
				return nil, fmt.Errorf("some-paths-builder-error")
			},
			expect:          func(fr *MockFileRepository) {},
			wantErrContains: "some-paths-builder-error",
		},
		{
			name: "errors if file exists check fails",
			givenConfig: game_launcher.Config{
				InstallPath:    "C:\\Program Files\\Publisher\\Game",
				ExecutableName: "Game.exe",
			},
			givenPathsBuilder: func(config game_launcher.Config) ([]string, error) {
				return []string{filepath.Join(config.InstallPath, "Game.running")}, nil
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().FileExists(gomock.Eq("C:\\Program Files\\Publisher\\Game\\Game.running")).Return(false, fmt.Errorf("some-io-error"))
			},
			wantErrContains: "some-io-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			mockRepository := NewMockFileRepository(ctrl)
			u := &url.URL{Host: net.JoinHostPort("1.1.1.1", "28960")}
			handler := MakeDeleteFileHookHandler(tt.givenPathsBuilder)
			args := map[string]string{}

			// EXPECT
			tt.expect(mockRepository)

			// WHEN
			err := handler.Run(mockRepository, u, tt.givenConfig, game_launcher.LaunchTypeLaunchAndJoin, args)

			if tt.wantErrContains != "" {
				assert.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
