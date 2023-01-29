//go:build unit

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			originURL := buildOriginURL(tt.givenOfferIDs, tt.givenArgs)
			assert.Equal(t, tt.expectedURL, originURL)
		})
	}
}
