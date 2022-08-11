//go:build unit

package internal

import (
	"fmt"
	"testing"

	"github.com/cetteup/joinme.click-launcher/pkg/refractor_config_handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEncryptedProfileConLogin(t *testing.T) {
	type test struct {
		name                 string
		prepareProfileConMap func(profileCon *refractor_config_handler.Config)
		wantErrContains      string
	}

	tests := []test{
		{
			name:                 "successfully extracts encrypted login details",
			prepareProfileConMap: func(profileCon *refractor_config_handler.Config) {},
		},
		{
			name: "fails if nickname is missing",
			prepareProfileConMap: func(profileCon *refractor_config_handler.Config) {
				profileCon.Delete(profileConKeyGamespyNick)
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if nickname is empty",
			prepareProfileConMap: func(profileCon *refractor_config_handler.Config) {
				profileCon.SetValue(profileConKeyGamespyNick, *refractor_config_handler.NewValue(""))
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if password is missing",
			prepareProfileConMap: func(profileCon *refractor_config_handler.Config) {
				profileCon.Delete(profileConKeyPassword)
			},
			wantErrContains: "encrypted password is missing/empty",
		},
		{
			name: "fails if password is empty",
			prepareProfileConMap: func(profileCon *refractor_config_handler.Config) {
				profileCon.SetValue(profileConKeyPassword, *refractor_config_handler.NewValue(""))
			},
			wantErrContains: "encrypted password is missing/empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			bytes := []byte(fmt.Sprintf("%s \"mister249\"\r\n%s \"some-encrypted-password\"\r\n", profileConKeyGamespyNick, profileConKeyPassword))
			profileCon := refractor_config_handler.ConfigFromBytes(bytes)
			tt.prepareProfileConMap(profileCon)

			// WHEN
			nickname, encryptedPassword, err := GetEncryptedProfileConLogin(profileCon)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				expectedNickname, err := profileCon.GetValue(profileConKeyGamespyNick)
				require.NoError(t, err)
				assert.Equal(t, expectedNickname.String(), nickname)
				expectedPassword, err := profileCon.GetValue(profileConKeyPassword)
				require.NoError(t, err)
				assert.Equal(t, expectedPassword.String(), encryptedPassword)
			}
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
