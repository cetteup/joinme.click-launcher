//go:build unit

package refractor_config_handler

import (
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

func TestHandler_ReadGlobalConfig(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		expect          func(repository *MockfileRepository, documentsDirPath string)
		wantErrContains string
	}

	tests := []test{
		{
			name:      "successfully reads config file",
			givenGame: GameBf2,
			expect: func(repository *MockfileRepository, documentsDirPath string) {
				repository.EXPECT().ReadFile(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, globalConFileName))).Return([]byte{}, nil)
			},
		},
		{
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			expect:          func(repository *MockfileRepository, documentsDirPath string) {},
			wantErrContains: "game not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(mockRepository, documentsDirPath)

			// WHEN
			config, err := handler.ReadGlobalConfig(tt.givenGame)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, map[string]Value{}, config.content)
			}
		})
	}
}

func TestHandler_ReadProfileConfig(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		givenProfile    string
		expect          func(repository *MockfileRepository, documentsDirPath string)
		wantErrContains string
	}

	tests := []test{
		{
			name:         "successfully reads config file",
			givenGame:    GameBf2,
			givenProfile: "0001",
			expect: func(repository *MockfileRepository, documentsDirPath string) {
				repository.EXPECT().ReadFile(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return([]byte{}, nil)
			},
		},
		{
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			expect:          func(repository *MockfileRepository, documentsDirPath string) {},
			wantErrContains: "game not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(mockRepository, documentsDirPath)

			// WHEN
			config, err := handler.ReadProfileConfig(tt.givenGame, tt.givenProfile)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, map[string]Value{}, config.content)
			}
		})
	}
}

func getHandlerWithDependencies(t *testing.T) (*Handler, *MockfileRepository) {
	ctrl := gomock.NewController(t)
	mockRepository := NewMockfileRepository(ctrl)
	return New(mockRepository), mockRepository
}
