package titles

import (
	"fmt"
	"net"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cetteup/joinme.click-launcher/internal/testhelpers"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
)

func TestBf2SetDefaultProfileHookHandler(t *testing.T) {
	type test struct {
		name            string
		givenArgs       map[string]string
		expect          func(fr *MockFileRepository)
		wantErrContains string
	}

	tests := []test{
		{
			name: "sets given profile as default profile",
			givenArgs: map[string]string{
				"profile": "0001",
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con")).Return([]byte("GlobalSettings.setDefaultUser \"0002\""), nil)
				fr.EXPECT().WriteFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con"), []byte("GlobalSettings.setDefaultUser \"0001\"\r\n"), gomock.Any())
			},
		},
		{
			name:            "errors if profile argument is missing",
			givenArgs:       map[string]string{},
			expect:          func(fr *MockFileRepository) {},
			wantErrContains: "required argument profile for hook set-default-profile is missing",
		},
		{
			name: "errors if Global.con cannot be read",
			givenArgs: map[string]string{
				"profile": "0001",
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con")).Return(nil, fmt.Errorf("some-read-error"))
			},
			wantErrContains: "some-read-error",
		},
		{
			name: "errors if Global.con cannot be written",
			givenArgs: map[string]string{
				"profile": "0001",
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con")).Return([]byte("GlobalSettings.setDefaultUser \"0002\""), nil)
				fr.EXPECT().WriteFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con"), []byte("GlobalSettings.setDefaultUser \"0001\"\r\n"), gomock.Any()).Return(fmt.Errorf("some-write-error"))
			},
			wantErrContains: "some-write-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			mockRepository := NewMockFileRepository(ctrl)
			u := &url.URL{Host: net.JoinHostPort("1.1.1.1", "16567")}
			config := game_launcher.Config{}
			handler := bf2SetDefaultProfileHookHandler{}

			// EXPECT
			tt.expect(mockRepository)

			// WHEN
			err := handler.Run(mockRepository, u, config, game_launcher.LaunchTypeLaunchAndJoin, tt.givenArgs)

			// THEN
			if tt.wantErrContains != "" {
				assert.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBf2PurgeServerHistoryHookHandler(t *testing.T) {
	type test struct {
		name            string
		givenArgs       map[string]string
		expect          func(fr *MockFileRepository)
		wantErrContains string
	}

	tests := []test{
		{
			name: "purges server history for given profile",
			givenArgs: map[string]string{
				"profile": "0001",
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0001\\General.con")).Return([]byte{}, nil)
				fr.EXPECT().WriteFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0001\\General.con"), gomock.Any(), gomock.Any())
			},
		},
		{
			name:      "purges server history for default profile",
			givenArgs: map[string]string{},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con")).Return([]byte("GlobalSettings.setDefaultUser \"0002\""), nil)
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0002\\General.con")).Return([]byte{}, nil)
				fr.EXPECT().WriteFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0002\\General.con"), gomock.Any(), gomock.Any())
			},
		},
		{
			name: "errors if given profile's General.con cannot be read",
			givenArgs: map[string]string{
				"profile": "0001",
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0001\\General.con")).Return(nil, fmt.Errorf("some-read-error"))
			},
			wantErrContains: "some-read-error",
		},
		{
			name:      "errors if default profile's General.con cannot be read",
			givenArgs: map[string]string{},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con")).Return([]byte("GlobalSettings.setDefaultUser \"0002\""), nil)
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0002\\General.con")).Return(nil, fmt.Errorf("some-read-error"))
			},
			wantErrContains: "some-read-error",
		},
		{
			name: "errors if profile's General.con cannot be written",
			givenArgs: map[string]string{
				"profile": "0001",
			},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0001\\General.con")).Return([]byte{}, nil)
				fr.EXPECT().WriteFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\0001\\General.con"), gomock.Any(), gomock.Any()).Return(fmt.Errorf("some-write-error"))
			},
			wantErrContains: "some-write-error",
		},
		{
			name:      "errors if Global.con cannot be read",
			givenArgs: map[string]string{},
			expect: func(fr *MockFileRepository) {
				fr.EXPECT().ReadFile(testhelpers.StringContainsMatcher("Battlefield 2\\Profiles\\Global.con")).Return(nil, fmt.Errorf("some-read-error"))
			},
			wantErrContains: "some-read-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			mockRepository := NewMockFileRepository(ctrl)
			u := &url.URL{Host: net.JoinHostPort("1.1.1.1", "16567")}
			config := game_launcher.Config{}
			handler := bf2PurgeServerHistoryHookHandler{}

			// EXPECT
			tt.expect(mockRepository)

			// WHEN
			err := handler.Run(mockRepository, u, config, game_launcher.LaunchTypeLaunchAndJoin, tt.givenArgs)

			// THEN
			if tt.wantErrContains != "" {
				assert.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
