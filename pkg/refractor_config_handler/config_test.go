package refractor_config_handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFromBytes(t *testing.T) {
	type test struct {
		name           string
		givenData      string
		expectedConfig Config
	}

	tests := []test{
		{
			name:      "parses config with unix line breaks",
			givenData: "GlobalSettings.setDefaultUser \"0010\"\nGlobalSettings.setNamePrefix \"=PRE=\"\n",
			expectedConfig: Config{
				content: map[string]Value{
					"GlobalSettings.setDefaultUser": {content: "\"0010\""},
					"GlobalSettings.setNamePrefix":  {content: "\"=PRE=\""},
				},
			},
		},
		{
			name:      "parses config with windows line breaks",
			givenData: "GlobalSettings.setDefaultUser \"0010\"\r\nGlobalSettings.setNamePrefix \"=PRE=\"\r\n",
			expectedConfig: Config{
				content: map[string]Value{
					"GlobalSettings.setDefaultUser": {content: "\"0010\""},
					"GlobalSettings.setNamePrefix":  {content: "\"=PRE=\""},
				},
			},
		},
		{
			name:      "parses multiple lines with same key",
			givenData: "GeneralSettings.setPlayedVOHelp \"HUD_HELP_A\"\nGeneralSettings.setPlayedVOHelp \"HUD_HELP_B\"\n",
			expectedConfig: Config{
				content: map[string]Value{
					"GeneralSettings.setPlayedVOHelp": {content: "\"HUD_HELP_A\";\"HUD_HELP_B\""},
				},
			},
		},
		{
			name:      "parses empty config",
			givenData: "",
			expectedConfig: Config{
				content: map[string]Value{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			config := ConfigFromBytes([]byte(tt.givenData))

			// THEN
			assert.Equal(t, &tt.expectedConfig, config)
		})
	}
}

func TestConfig_GetValue(t *testing.T) {
	type test struct {
		name            string
		givenConfig     Config
		givenKey        string
		expectedValue   Value
		wantErrContains string
	}

	tests := []test{
		{
			name: "successfully retrieves value",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey:      "some-key",
			expectedValue: Value{content: "some-value"},
		},
		{
			name: "error for non-existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey:        "some-other-key",
			wantErrContains: "no such key in config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			value, err := tt.givenConfig.GetValue(tt.givenKey)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}

func TestConfig_SetValue(t *testing.T) {
	type test struct {
		name           string
		givenConfig    Config
		givenKey       string
		givenValue     Value
		expectedConfig Config
	}

	tests := []test{
		{
			name: "adds value under new key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey:   "other-key",
			givenValue: Value{content: "other-value"},
			expectedConfig: Config{
				content: map[string]Value{
					"some-key":  {content: "some-value"},
					"other-key": {content: "other-value"},
				},
			},
		},
		{
			name: "overwrites value at existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "old-value"},
				},
			},
			givenKey:   "some-key",
			givenValue: Value{content: "new-value"},
			expectedConfig: Config{
				content: map[string]Value{
					"some-key": {content: "new-value"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			config := tt.givenConfig

			// WHEN
			config.SetValue(tt.givenKey, tt.givenValue)

			// THEN
			assert.Equal(t, tt.expectedConfig, config)
		})
	}
}

func TestConfig_Delete(t *testing.T) {
	type test struct {
		name           string
		givenConfig    Config
		givenKey       string
		expectedConfig Config
	}

	tests := []test{
		{
			name: "removes existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey: "some-key",
			expectedConfig: Config{
				content: map[string]Value{},
			},
		},
		{
			name: "noop for non-existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey: "other-key",
			expectedConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			config := tt.givenConfig

			// WHEN
			config.Delete(tt.givenKey)

			// THEN
			assert.Equal(t, tt.expectedConfig, config)
		})
	}
}

func TestValue_String(t *testing.T) {
	type test struct {
		name           string
		givenValue     Value
		expectedString string
	}

	tests := []test{
		{
			name:           "returns non-quoted string as is",
			givenValue:     Value{content: "some-unquoted-value"},
			expectedString: "some-unquoted-value",
		},
		{
			name:           "returns string containing quotes as is",
			givenValue:     Value{content: "\"some-quoted-sub-value\" some-unquoted-sub-value"},
			expectedString: "\"some-quoted-sub-value\" some-unquoted-sub-value",
		},
		{
			name:           "returns quoted string without quotes",
			givenValue:     Value{content: "\"some-quoted-value\""},
			expectedString: "some-quoted-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			asString := tt.givenValue.String()

			// THEN
			assert.Equal(t, tt.expectedString, asString)
		})
	}
}

func TestValue_Slice(t *testing.T) {
	type test struct {
		name          string
		givenValue    Value
		expectedSlice []string
	}

	tests := []test{
		{
			name:          "returns unquoted single value string as is in slice with one element",
			givenValue:    Value{content: "some-unquoted-single-value"},
			expectedSlice: []string{"some-unquoted-single-value"},
		},
		{
			name:          "returns single value string containing quotes as is in slice with one element",
			givenValue:    Value{content: "\"some-quoted-single-sub-value\" some-unquoted-single-sub-value"},
			expectedSlice: []string{"\"some-quoted-single-sub-value\" some-unquoted-single-sub-value"},
		},
		{
			name:          "returns quoted single value string without quotes in slice with one element",
			givenValue:    Value{content: "\"some-quoted-single-value\""},
			expectedSlice: []string{"some-quoted-single-value"},
		},
		{
			name:          "returns unquoted multi value string as is in slice with multiple elements",
			givenValue:    Value{content: "some-unquoted-value;some-other-unquoted-value"},
			expectedSlice: []string{"some-unquoted-value", "some-other-unquoted-value"},
		},
		{
			name:          "returns multi value string containing quotes as is in slice with multiple elements",
			givenValue:    Value{content: "\"some-quoted-sub-value\" some-unquoted-sub-value;\"some-other-quoted-sub-value\" some-other-unquoted-sub-value"},
			expectedSlice: []string{"\"some-quoted-sub-value\" some-unquoted-sub-value", "\"some-other-quoted-sub-value\" some-other-unquoted-sub-value"},
		},
		{
			name:          "returns quoted multi value string without quotes in slice with multiple elements",
			givenValue:    Value{content: "\"some-quoted-value\";\"some-other-quoted-value\""},
			expectedSlice: []string{"some-quoted-value", "some-other-quoted-value"},
		},
		{
			name:          "returns mixed quoted multi value string without quotes and as is in slice with multiple elements",
			givenValue:    Value{content: "\"some-quoted-value\";some-unquoted-value"},
			expectedSlice: []string{"some-quoted-value", "some-unquoted-value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			asSlice := tt.givenValue.Slice()

			// THEN
			assert.Equal(t, tt.expectedSlice, asSlice)
		})
	}
}
