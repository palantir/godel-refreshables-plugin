// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package config

import (
	"io/ioutil"

	v0 "github.com/palantir/godel-refreshables-plugin/config/internal/v0"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config v0.Config

func ToConfig(in *Config) *v0.Config {
	return (*v0.Config)(in)
}

func ReadConfigFromFile(f string) (Config, error) {
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return Config{}, errors.WithStack(err)
	}
	return ReadConfigFromBytes(bytes)
}

func ReadConfigFromBytes(inputBytes []byte) (Config, error) {
	var cfg Config
	if err := yaml.UnmarshalStrict(inputBytes, &cfg); err != nil {
		return Config{}, errors.WithStack(err)
	}
	return cfg, nil
}
