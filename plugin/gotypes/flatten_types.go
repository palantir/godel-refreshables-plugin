// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gotypes

import (
	"go/types"

	"github.com/palantir/godel-refreshables-plugin/plugin/tags"
)

// FlattenTypes recursively iterates through type's fields and underlying types and returns a list.
func FlattenTypes(types ...types.Type) []types.Type {
	cont := typesContainer{}
	for _, typ := range types {
		cont.flattenTypesRecursive(typ)
	}
	return cont.types
}

type typesContainer struct {
	types []types.Type
}

func (t *typesContainer) flattenTypesRecursive(typ types.Type) {
	for _, seenType := range t.types {
		if types.Identical(seenType, typ) {
			// return early if we've seen this type already
			return
		}
	}

	typeStr := typ.String()
	typeUnderlying := typ.Underlying()
	_, _ = typeStr, typeUnderlying

	switch underlying := typ.(type) {
	case *types.Basic:
		t.types = append(t.types, typ)
	case *types.Named:
		t.types = append(t.types, typ)
		t.flattenTypesRecursive(underlying.Underlying())
	case *types.Pointer:
		t.types = append(t.types, typ)
		t.flattenTypesRecursive(underlying.Elem())
	case *types.Map:
		t.types = append(t.types, typ)
		t.flattenTypesRecursive(underlying.Key())
		t.flattenTypesRecursive(underlying.Elem())
	case *types.Array:
		t.types = append(t.types, typ)
		t.flattenTypesRecursive(underlying.Elem())
	case *types.Slice:
		t.types = append(t.types, typ)
		t.flattenTypesRecursive(underlying.Elem())
	case *types.Struct:
		for i := 0; i < underlying.NumFields(); i++ {
			fieldOptions, _ := tags.ParseTag(underlying.Tag(i))
			if fieldOptions.Exclude {
				continue
			}
			t.flattenTypesRecursive(underlying.Field(i).Type())
		}
	}
}
