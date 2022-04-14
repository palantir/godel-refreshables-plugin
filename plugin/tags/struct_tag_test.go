// Copyright (c) 2022 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	for _, currCase := range []struct {
		name        string
		tagIn       string
		expectedOut RefreshablesFieldTagOptions
		expectedErr string
	}{
		{
			name:        "empty",
			tagIn:       "",
			expectedOut: RefreshablesFieldTagOptions{},
		},
		{
			name:        "only other tags",
			tagIn:       `yaml:"foo" json:"bar"`,
			expectedOut: RefreshablesFieldTagOptions{},
		},
		{
			name:  "options with name",
			tagIn: `refreshables:"custom-name"`,
			expectedOut: RefreshablesFieldTagOptions{
				Name: "custom-name",
			},
		},
		{
			name:  "excluded without name",
			tagIn: `refreshables:",exclude"`,
			expectedOut: RefreshablesFieldTagOptions{
				Exclude: true,
			},
		},
		{
			name:  "exclude with name and yaml tag",
			tagIn: `yaml:"foo" refreshables:"bar,exclude"`,
			expectedOut: RefreshablesFieldTagOptions{
				Name:    "bar",
				Exclude: true,
			},
		},
		{
			name:        "invalid refreshable option",
			tagIn:       `refreshables:"custom,invalidoption"`,
			expectedErr: "invalid tag option: \"invalidoption\" in tag `refreshables:\"custom,invalidoption\"`",
		},
	} {
		t.Run(currCase.name, func(t *testing.T) {
			out, err := ParseTag(currCase.tagIn)
			if currCase.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, currCase.expectedOut, out)
			} else {
				assert.Error(t, err, currCase.expectedErr)
			}
		})
	}
}
