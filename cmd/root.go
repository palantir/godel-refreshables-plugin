// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/palantir/godel/v2/framework/pluginapi"
	"github.com/palantir/pkg/cobracli"
	"github.com/spf13/cobra"
)

var (
	projectDirFlagVal string
	configFileFlagVal string
)

var rootCmd = &cobra.Command{
	Use: "refreshables-plugin",
}

func Execute() int {
	return cobracli.ExecuteWithDefaultParams(rootCmd)
}

func init() {
	pluginapi.AddProjectDirPFlagPtr(rootCmd.PersistentFlags(), &projectDirFlagVal)
	if err := rootCmd.MarkPersistentFlagRequired(pluginapi.ProjectDirFlagName); err != nil {
		panic(err)
	}
	pluginapi.AddConfigPFlagPtr(rootCmd.PersistentFlags(), &configFileFlagVal)
	if err := rootCmd.MarkPersistentFlagRequired(pluginapi.ConfigFlagName); err != nil {
		panic(err)
	}
}
