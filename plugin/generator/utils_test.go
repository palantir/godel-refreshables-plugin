// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package generator

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongestCommonPkgPathSuffix(t *testing.T) {
	for _, test := range []struct {
		name                  string
		testPkgA              []string
		testPkgB              []string
		expectedLongestSuffix int
	}{
		{
			name:                  "nothing in common returns 0",
			testPkgA:              strings.Split("encoding/json", string(os.PathSeparator)),
			testPkgB:              strings.Split("bytes", string(os.PathSeparator)),
			expectedLongestSuffix: 0,
		},
		{
			name:                  "common last component returns 1",
			testPkgA:              strings.Split("github.com/abradshaw/go-refreshable/generator/testtypes/anotherpackage/config", string(os.PathSeparator)),
			testPkgB:              strings.Split("github.com/abradshaw/go-refreshable/generator/testtypes/config", string(os.PathSeparator)),
			expectedLongestSuffix: 1,
		},
		{
			name:                  "handle mismatched input lengths",
			testPkgA:              strings.Split("testtypes/anotherpackage/config", string(os.PathSeparator)),
			testPkgB:              strings.Split("config", string(os.PathSeparator)),
			expectedLongestSuffix: 1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actualLongestSuffix := longestCommonPkgPathSuffix(test.testPkgA, test.testPkgB)
			assert.Equal(t, test.expectedLongestSuffix, actualLongestSuffix)
		})
	}
}
