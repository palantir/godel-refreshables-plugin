// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package config_test

import (
	"testing"

	"github.com/palantir/godel-refreshables-plugin/config"
	v0 "github.com/palantir/godel-refreshables-plugin/config/internal/v0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestReadConfig(t *testing.T) {
	for i, tc := range []struct {
		in   string
		want config.Config
	}{
		{
			`
refreshables:
  ./config:
    types:
      - RuntimeConfig
`,
			config.Config{
				Refreshables: map[string]v0.RefreshablePackageConfig{
					"./config": {
						Types:  []string{"RuntimeConfig"},
						Output: "",
					},
				},
			},
		},
		{
			`
refreshables:
  ./config:
    types:
      - RuntimeConfig
    output: ./generated_src
`,
			config.Config{
				Refreshables: map[string]v0.RefreshablePackageConfig{
					"./config": {
						Types:  []string{"RuntimeConfig"},
						Output: "./generated_src",
					},
				},
			},
		},
	} {
		var got config.Config
		err := yaml.Unmarshal([]byte(tc.in), &got)
		require.NoError(t, err)
		assert.Equal(t, tc.want, got, "Case %d", i)
	}
}
