// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/palantir/godel-refreshables-plugin/config"
	"github.com/palantir/godel-refreshables-plugin/plugin"
	"github.com/spf13/cobra"
)

const (
	generateCmdName = "generate"
	verifyFlagName  = "verify"
)

var (
	verifyFlag bool
)

var generateCmd = &cobra.Command{
	Use:   generateCmdName,
	Short: "Generates refreshable types for structs listed in plugin configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.ReadConfigFromFile(configFileFlagVal)
		if err != nil {
			return err
		}
		if len(cfg.Refreshables) == 0 {
			return nil
		}
		return plugin.Run(projectDirFlagVal, cfg, verifyFlag)
	},
}

func init() {
	generateCmd.Flags().BoolVar(&verifyFlag, verifyFlagName, false, "verify that current project matches output of the generator")
	rootCmd.AddCommand(generateCmd)
}
