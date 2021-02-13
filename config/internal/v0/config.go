// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package v0

type Config struct {
	Refreshables map[string]RefreshablePackageConfig `yaml:"refreshables,omitempty"`
}

type RefreshablePackageConfig struct {
	// Types is a list of type names within the package for which refreshables will be generated.
	// All required element types will be generated as needed.
	Types []string `yaml:"types,omitempty"`
	// Output is an optional override for the output file path to write the generated code.
	// By default, the generated refreshables will be written within the package as zz_generated_refreshables.go.
	Output string `yaml:"output,omitempty"`
}
