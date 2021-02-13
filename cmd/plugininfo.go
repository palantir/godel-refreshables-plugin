// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/palantir/godel/v2/framework/pluginapi/v2/pluginapi"
	"github.com/palantir/godel/v2/framework/verifyorder"
)

const (
	pluginDescription = "Generates interfaces and implementations for strongly-typing refreshable objects."
)

var (
	Version    = "unspecified"
	PluginInfo = pluginapi.MustNewPluginInfo(
		"com.palantir.godel-refreshables-plugin",
		"refreshables-plugin",
		Version,
		pluginapi.PluginInfoUsesConfigFile(),
		pluginapi.PluginInfoGlobalFlagOptions(
			pluginapi.GlobalFlagOptionsParamDebugFlag("--"+pluginapi.DebugFlagName),
			pluginapi.GlobalFlagOptionsParamProjectDirFlag("--"+pluginapi.ProjectDirFlagName),
			pluginapi.GlobalFlagOptionsParamConfigFlag("--"+pluginapi.ConfigFlagName),
		),
		pluginapi.PluginInfoTaskInfo(
			"refreshables",
			pluginDescription,
			pluginapi.TaskInfoCommand("generate"),
			pluginapi.TaskInfoVerifyOptions(
				pluginapi.VerifyOptionsOrdering(pInt(verifyorder.Generate+1)),
				pluginapi.VerifyOptionsApplyFalseArgs("--"+verifyFlagName),
			),
		),
	)
)

func pInt(i int) *int { return &i }
