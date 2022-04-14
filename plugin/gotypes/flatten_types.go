// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gotypes

import (
	"go/types"

	"github.com/palantir/godel-refreshables-plugin/plugin/tags"
)

// FlattenTypes recursively iterates through type's fields and underlying types and returns a list.
func FlattenTypes(types ...types.Type) ([]types.Type, error) {
	cont := typesContainer{}
	for _, typ := range types {
		if err := cont.flattenTypesRecursive(typ); err != nil {
			return nil, err
		}
	}
	return cont.types, nil
}

type typesContainer struct {
	types []types.Type
}

func (t *typesContainer) flattenTypesRecursive(typ types.Type) error {
	for _, seenType := range t.types {
		if types.Identical(seenType, typ) {
			// return early if we've seen this type already
			return nil
		}
	}

	switch underlying := typ.(type) {
	case *types.Basic:
		t.types = append(t.types, typ)
	case *types.Named:
		t.types = append(t.types, typ)
		if err := t.flattenTypesRecursive(underlying.Underlying()); err != nil {
			return err
		}
	case *types.Pointer:
		t.types = append(t.types, typ)
		if err := t.flattenTypesRecursive(underlying.Elem()); err != nil {
			return err
		}
	case *types.Map:
		t.types = append(t.types, typ)
		if err := t.flattenTypesRecursive(underlying.Key()); err != nil {
			return err
		}
		if err := t.flattenTypesRecursive(underlying.Elem()); err != nil {
			return err
		}
	case *types.Array:
		t.types = append(t.types, typ)
		if err := t.flattenTypesRecursive(underlying.Elem()); err != nil {
			return err
		}
	case *types.Slice:
		t.types = append(t.types, typ)
		if err := t.flattenTypesRecursive(underlying.Elem()); err != nil {
			return err
		}
	case *types.Struct:
		for i := 0; i < underlying.NumFields(); i++ {
			fieldOptions, err := tags.ParseTag(underlying.Tag(i))
			if err != nil {
				return err
			}
			if fieldOptions.Exclude {
				continue
			}
			if err := t.flattenTypesRecursive(underlying.Field(i).Type()); err != nil {
				return err
			}
		}
	}
	return nil
}
