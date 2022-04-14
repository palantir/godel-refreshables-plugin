// Copyright (c) 2022 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package tags

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// RefreshablesFieldTagOptions represents custom struct options to change behavior of refreshables code generation.
type RefreshablesFieldTagOptions struct {
	// Name is currently not used, but this might be used in the future to override the name of generated methods.
	Name string
	// Exclude is used to indicate that the field should be excluded from generated code.
	Exclude bool
}

// ParseTag returns the RefreshablesFieldTagOptions from the tag.
// This uses a field tag syntax similar to json or yaml struct tags.
//
//   type Example struct {
//   	Foo FooType `yaml:"foo,omitempty" refreshables:",exclude"`
//   }
func ParseTag(tag string) (RefreshablesFieldTagOptions, error) {
	result := RefreshablesFieldTagOptions{}
	structTag := reflect.StructTag(tag)
	if value, ok := structTag.Lookup("refreshables"); ok {
		parts := strings.Split(value, ",")
		result.Name = parts[0]
		for _, option := range parts[1:] {
			if option == "" {
				// ignore
			} else if option == "exclude" {
				result.Exclude = true
			} else {
				return result, errors.Errorf("invalid tag option: \"%s\" in tag `%s`", option, tag)
			}
		}
	}

	return result, nil
}
