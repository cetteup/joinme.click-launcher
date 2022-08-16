//go:build unit

package software_finder

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

func TestSoftwareFinder_IsInstalledAnywhere(t *testing.T) {
	type test struct {
		name                    string
		givenConfigs            []Config
		expect                  func(rr *MockRegistryRepository, fr *MockFileRepository)
		wantIsInstalledAnywhere bool
		wantErrContains         string
	}

	tests := []test{
		{
			name: "true for installed software with single config",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("C:\\Some\\Game", nil)
			},
			wantIsInstalledAnywhere: true,
		},
		{
			name: "true for installed software with multiple configs",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			wantIsInstalledAnywhere: true,
		},
		{
			name: "false if no config matches",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, nil)
			},
			wantIsInstalledAnywhere: false,
		},
		{
			name: "silently errors if there are more configs",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", fmt.Errorf("some-error-that-is-not-returned"))
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			wantIsInstalledAnywhere: true,
		},
		{
			name: "errors if last config errors",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, fmt.Errorf("some-error-that-is-returned"))
			},
			wantErrContains: "some-error-that-is-returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			finder, mockRegistryRepository, mockFileRepository := getFinderWithDependencies(t)

			// EXPECT
			tt.expect(mockRegistryRepository, mockFileRepository)

			// WHEN
			isInstalledAnywhere, err := finder.IsInstalledAnywhere(tt.givenConfigs)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantIsInstalledAnywhere, isInstalledAnywhere)
			}
		})
	}
}

func TestSoftwareFinder_IsInstalled(t *testing.T) {
	type test struct {
		name            string
		givenConfig     Config
		expect          func(rr *MockRegistryRepository, fr *MockFileRepository)
		wantIsInstalled bool
		wantErrContains string
	}

	tests := []test{
		{
			name: "true for installed software via registry finder",
			givenConfig: Config{
				ForType:           RegistryFinder,
				RegistryPath:      "SOFTWARE\\some\\game",
				RegistryValueName: "InstallDir",
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("C:\\Some\\Game", nil)
			},
			wantIsInstalled: true,
		},
		{
			name: "true for installed software via path finder using directory",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game",
				PathType:    PathTypeDir,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			wantIsInstalled: true,
		},
		{
			name: "true for installed software via path finder using file",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game\\launch.exe",
				PathType:    PathTypeFile,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().FileExists("C:\\Some\\Game\\launch.exe").Return(true, nil)
			},
			wantIsInstalled: true,
		},
		{
			name: "false for non-installed software via registry finder",
			givenConfig: Config{
				ForType:           RegistryFinder,
				RegistryPath:      "SOFTWARE\\some\\game",
				RegistryValueName: "InstallDir",
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
			},
			wantIsInstalled: false,
		},
		{
			name: "false for non-installed software via path finder using directory",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game",
				PathType:    PathTypeDir,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, nil)
			},
			wantIsInstalled: false,
		},
		{
			name: "false for non-installed software via path finder using file",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game\\launch.exe",
				PathType:    PathTypeFile,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().FileExists("C:\\Some\\Game\\launch.exe").Return(false, nil)
			},
			wantIsInstalled: false,
		},
		{
			name: "errors for registry error",
			givenConfig: Config{
				ForType:           RegistryFinder,
				RegistryPath:      "SOFTWARE\\some\\game",
				RegistryValueName: "InstallDir",
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name: "errors for path finder error using directory",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game",
				PathType:    PathTypeDir,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name: "errors for path finder error using file",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game\\launch.exe",
				PathType:    PathTypeFile,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().FileExists("C:\\Some\\Game\\launch.exe").Return(false, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name: "errors for unsupported path type",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game\\launch.exe",
				PathType:    -1,
			},
			expect:          func(rr *MockRegistryRepository, fr *MockFileRepository) {},
			wantErrContains: "unsupported path type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			finder, mockRegistryRepository, mockFileRepository := getFinderWithDependencies(t)

			// EXPECT
			tt.expect(mockRegistryRepository, mockFileRepository)

			// WHEN
			isInstalled, err := finder.IsInstalled(tt.givenConfig)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantIsInstalled, isInstalled)
			}
		})
	}
}

func TestSoftwareFinder_GetInstallDirFromSomewhere(t *testing.T) {
	type test struct {
		name               string
		givenConfigs       []Config
		expect             func(rr *MockRegistryRepository, fr *MockFileRepository)
		expectedInstallDir string
		wantErrContains    string
	}

	tests := []test{
		{
			name: "successfully determines install dir with single config",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("C:\\Some\\Game", nil)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			expectedInstallDir: "C:\\Some\\Game",
		},
		{
			name: "successfully determines install dir with multiple configs",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			expectedInstallDir: "C:\\Some\\Game",
		},
		{
			name: "silently errors if there are more configs",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", fmt.Errorf("some-error-that-is-not-returned"))
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			expectedInstallDir: "C:\\Some\\Game",
		},
		{
			name: "errors if there are no more configs",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, fmt.Errorf("some-error-that-is-returned"))
			},
			wantErrContains: "some-error-that-is-returned",
		},
		{
			name: "errors if no configs matches",
			givenConfigs: []Config{
				{
					ForType:           RegistryFinder,
					RegistryPath:      "SOFTWARE\\some\\game",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     PathFinder,
					InstallPath: "C:\\Some\\Game",
					PathType:    PathTypeDir,
				},
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, nil)
				fr.EXPECT().DirExists("C:\\Some").Return(false, nil)
			},
			wantErrContains: "failed to determine install path based on received path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			finder, mockRegistryRepository, mockFileRepository := getFinderWithDependencies(t)

			// EXPECT
			tt.expect(mockRegistryRepository, mockFileRepository)

			// WHEN
			installDir, err := finder.GetInstallDirFromSomewhere(tt.givenConfigs)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedInstallDir, installDir)
			}
		})
	}
}

func TestSoftwareFinder_GetInstallDir(t *testing.T) {
	type test struct {
		name               string
		givenConfig        Config
		expect             func(rr *MockRegistryRepository, fr *MockFileRepository)
		expectedInstallDir string
		wantErrContains    string
	}

	tests := []test{
		{
			name: "successfully determines install dir via registry finder",
			givenConfig: Config{
				ForType:           RegistryFinder,
				RegistryPath:      "SOFTWARE\\some\\game",
				RegistryValueName: "InstallDir",
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("C:\\Some\\Game", nil)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			expectedInstallDir: "C:\\Some\\Game",
		},
		{
			name: "successfully determines install dir via registry finder with file path value",
			givenConfig: Config{
				ForType:           RegistryFinder,
				RegistryPath:      "SOFTWARE\\some\\game",
				RegistryValueName: "InstallDir",
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("C:\\Some\\Game\\launch.exe", nil)
				fr.EXPECT().DirExists("C:\\Some\\Game\\launch.exe").Return(false, nil)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			expectedInstallDir: "C:\\Some\\Game",
		},
		{
			name: "successfully determines install dir via path finder using directory",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game",
				PathType:    PathTypeDir,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			expectedInstallDir: "C:\\Some\\Game",
		},
		{
			name: "successfully determines install dir via path finder using file",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game\\launch.exe",
				PathType:    PathTypeFile,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(true, nil)
			},
			expectedInstallDir: "C:\\Some\\Game",
		},
		{
			name: "errors for registry error",
			givenConfig: Config{
				ForType:           RegistryFinder,
				RegistryPath:      "SOFTWARE\\some\\game",
				RegistryValueName: "InstallDir",
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("", registry.ErrNotExist)
			},
			wantErrContains: "The system cannot find the file specified",
		},
		{
			name: "errors for path validation error",
			givenConfig: Config{
				ForType:           RegistryFinder,
				RegistryPath:      "SOFTWARE\\some\\game",
				RegistryValueName: "InstallDir",
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				rr.EXPECT().GetStringValue(registry.LOCAL_MACHINE, "SOFTWARE\\some\\game", "InstallDir").Return("C:\\Some\\Game", nil)
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name: "errors for path finder error",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game\\launch.exe",
				PathType:    PathTypeFile,
			},
			expect: func(rr *MockRegistryRepository, fr *MockFileRepository) {
				fr.EXPECT().DirExists("C:\\Some\\Game").Return(false, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name: "errors unsupported path type",
			givenConfig: Config{
				ForType:     PathFinder,
				InstallPath: "C:\\Some\\Game\\launch.exe",
				PathType:    -1,
			},
			expect:          func(rr *MockRegistryRepository, fr *MockFileRepository) {},
			wantErrContains: "unsupported path type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			finder, mockRegistryRepository, mockFileRepository := getFinderWithDependencies(t)

			// EXPECT
			tt.expect(mockRegistryRepository, mockFileRepository)

			// WHEN
			installDir, err := finder.GetInstallDir(tt.givenConfig)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedInstallDir, installDir)
			}
		})
	}
}

func getFinderWithDependencies(t *testing.T) (*SoftwareFinder, *MockRegistryRepository, *MockFileRepository) {
	ctrl := gomock.NewController(t)
	mockRegistryRepository := NewMockRegistryRepository(ctrl)
	mockFileRepository := NewMockFileRepository(ctrl)
	return New(mockRegistryRepository, mockFileRepository), mockRegistryRepository, mockFileRepository
}
