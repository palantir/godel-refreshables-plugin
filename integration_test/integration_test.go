// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package integration_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/palantir/godel-refreshables-plugin/config"
	"github.com/palantir/godel-refreshables-plugin/plugin"
	"github.com/palantir/godel/v2/pkg/products"
	"github.com/stretchr/testify/require"
)

func Test1(t *testing.T) {
	const (
		cfgFile = "testcode/test1/refreshables-plugin.yml"
	)
	cfg, err := config.ReadConfigFromFile(cfgFile)
	require.NoError(t, err)
	err = plugin.Run("./testcode/test1", cfg, false)
	require.NoError(t, err)

	cli, err := products.Bin("refreshables-plugin")
	require.NoError(t, err)

	cmd := exec.Command(cli, "generate", "--project-dir=./testcode/test1", "--config", cfgFile, "--verify")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "plugin verify failed: %s", string(out))
}

func TestHttpclient(t *testing.T) {
	const (
		cfgFile = "testcode/httpclient/refreshables-plugin.yml"
		outFile = "testcode/httpclient/zz_generated_refreshables.go"
	)
	// remove outFile if it exists
	_, err := os.Stat(outFile)
	if !os.IsNotExist(err) {
		_ = os.Remove(outFile)
	}

	cfg, err := config.ReadConfigFromFile(cfgFile)
	require.NoError(t, err)
	err = plugin.Run("./testcode/httpclient", cfg, false)
	require.NoError(t, err)
	_, err = os.Stat(outFile)
	require.NoError(t, err)

	cli, err := products.Bin("refreshables-plugin")
	require.NoError(t, err)

	cmd := exec.Command(cli, "generate", "--project-dir=./testcode/httpclient", "--config", cfgFile, "--verify")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "plugin verify failed: %s", string(out))
}
