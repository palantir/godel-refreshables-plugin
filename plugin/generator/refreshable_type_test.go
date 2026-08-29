// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package generator

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGeneratedFileHeaderMatchesGoConvention pins the header to the form that go, gopls and
// golangci-lint use to recognize a file as generated.
// See https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source.
func TestGeneratedFileHeaderMatchesGoConvention(t *testing.T) {
	file, err := GenerateRefreshableFile("github.com/palantir/example", "example", nil)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	require.NoError(t, file.Render(buf))

	header, _, _ := strings.Cut(buf.String(), "\n")
	assert.Regexp(t, `^// Code generated .* DO NOT EDIT\.$`, header)
}
