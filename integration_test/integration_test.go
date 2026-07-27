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
	"github.com/stretchr/testify/assert"
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

func Test2(t *testing.T) {
	const (
		cfgFile = "testcode/test2/refreshables-plugin.yml"
	)
	cfg, err := config.ReadConfigFromFile(cfgFile)
	require.NoError(t, err)
	err = plugin.Run("./testcode/test2", cfg, false)
	require.NoError(t, err)

	cli, err := products.Bin("refreshables-plugin")
	require.NoError(t, err)

	cmd := exec.Command(cli, "generate", "--project-dir=./testcode/test2", "--config", cfgFile, "--verify")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "plugin verify failed: %s", string(out))
}

func Test3(t *testing.T) {
	const (
		cfgFile = "testcode/test3/refreshables-plugin.yml"
	)
	cfg, err := config.ReadConfigFromFile(cfgFile)
	require.NoError(t, err)
	err = plugin.Run("./testcode/test3", cfg, false)
	require.NoError(t, err)

	cli, err := products.Bin("refreshables-plugin")
	require.NoError(t, err)

	cmd := exec.Command(cli, "generate", "--project-dir=./testcode/test3", "--config", cfgFile, "--verify")
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

func TestGenerateLeavesUpToDateFileAlone(t *testing.T) {
	const (
		cfgFile = "testcode/test1/refreshables-plugin.yml"
		outFile = "testcode/test1/zz_generated_refreshables.go"
	)
	restoreFile(t, outFile)
	appendLine(t, outFile, "// stale")

	cfg, err := config.ReadConfigFromFile(cfgFile)
	require.NoError(t, err)

	require.NoError(t, plugin.Run("./testcode/test1", cfg, false))
	first, err := os.Stat(outFile)
	require.NoError(t, err)

	require.NoError(t, plugin.Run("./testcode/test1", cfg, false))
	second, err := os.Stat(outFile)
	require.NoError(t, err)
	assert.Equal(t, first.ModTime(), second.ModTime(), "file whose contents already match was rewritten")
}

func TestVerifyReportsEveryOutdatedFileWithoutWriting(t *testing.T) {
	const cfgFile = "testcode/test1/refreshables-plugin.yml"
	outFiles := []string{
		"testcode/test1/zz_generated_refreshables.go",
		"testcode/test1/librarypkg/zz_generated_refreshables.go",
	}
	for _, outFile := range outFiles {
		restoreFile(t, outFile)
		appendLine(t, outFile, "// stale")
	}

	cfg, err := config.ReadConfigFromFile(cfgFile)
	require.NoError(t, err)
	err = plugin.Run("./testcode/test1", cfg, true)
	require.EqualError(t, err, `generated output is out of date:
librarypkg/zz_generated_refreshables.go: out of date
zz_generated_refreshables.go: out of date`)

	for _, outFile := range outFiles {
		contents, err := os.ReadFile(outFile)
		require.NoError(t, err)
		assert.Contains(t, string(contents), "// stale", "verify rewrote %s", outFile)
	}
}

// TestDuplicateOutputPathRejected covers two configured packages resolving to the same output file, which
// is a collision rather than a last writer that depends on map iteration order.
func TestDuplicateOutputPathRejected(t *testing.T) {
	const dupFile = "testcode/test3/zz_duplicate.go"
	t.Cleanup(func() {
		_ = os.Remove(dupFile)
	})

	cfg, err := config.ReadConfigFromBytes([]byte(`refreshables:
  .:
    output: ./zz_duplicate.go
    types:
      - MainStruct
  github.com/palantir/godel-refreshables-plugin/integration_test/testcode/test3:
    output: ./zz_duplicate.go
    types:
      - MainStruct
`))
	require.NoError(t, err)

	err = plugin.Run("./testcode/test3", cfg, false)
	require.EqualError(t, err, `generated path "zz_duplicate.go" was produced more than once`)
	assert.NoFileExists(t, dupFile)
}

// TestOutputEscapingProjectDirRejected uses an output path whose ".." segments are not leading ones, so
// only resolving the path catches that it lands outside the project directory.
func TestOutputEscapingProjectDirRejected(t *testing.T) {
	const escapedFile = "testcode/zz_escaped.go"
	t.Cleanup(func() {
		_ = os.Remove(escapedFile)
	})

	cfg, err := config.ReadConfigFromBytes([]byte(`refreshables:
  .:
    output: ./nested/../../zz_escaped.go
    types:
      - MainStruct
`))
	require.NoError(t, err)

	err = plugin.Run("./testcode/test3", cfg, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "escapes the project directory")
	assert.NoFileExists(t, escapedFile)
}

// restoreFile restores path's original contents when the test finishes, so a test that deliberately
// leaves generated output stale does not leak that state into the checked-in tree.
func restoreFile(t *testing.T, path string) {
	t.Helper()
	original, err := os.ReadFile(path)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.WriteFile(path, original, 0644))
	})
}

// appendLine makes a generated file differ from what the generator produces while keeping it valid Go,
// which the plugin needs in order to load the package it generates for.
func appendLine(t *testing.T, path, line string) {
	t.Helper()
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, append(contents, []byte(line+"\n")...), 0644))
}
