// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package generator

import (
	"go/types"
	"math"
	"strings"

	"github.com/dave/jennifer/jen"
)

const refreshablePath = "github.com/palantir/pkg/refreshable"

var (
	refreshable                       = jen.Qual(refreshablePath, "Refreshable")
	refreshableString                 = jen.Qual(refreshablePath, "String")
	refreshableStringConstructor      = jen.Qual(refreshablePath, "NewString")
	refreshableStringPtr              = jen.Qual(refreshablePath, "StringPtr")
	refreshableStringPtrConstructor   = jen.Qual(refreshablePath, "NewStringPtr")
	refreshableStringSlice            = jen.Qual(refreshablePath, "StringSlice")
	refreshableStringSliceConstructor = jen.Qual(refreshablePath, "NewStringSlice")
	refreshableBool                   = jen.Qual(refreshablePath, "Bool")
	refreshableBoolConstructor        = jen.Qual(refreshablePath, "NewBool")
	refreshableBoolPtr                = jen.Qual(refreshablePath, "BoolPtr")
	refreshableBoolPtrConstructor     = jen.Qual(refreshablePath, "NewBoolPtr")
	refreshableInt                    = jen.Qual(refreshablePath, "Int")
	refreshableIntConstructor         = jen.Qual(refreshablePath, "NewInt")
	refreshableIntPtr                 = jen.Qual(refreshablePath, "IntPtr")
	refreshableIntPtrConstructor      = jen.Qual(refreshablePath, "NewIntPtr")
	refreshableInt64                  = jen.Qual(refreshablePath, "Int64")
	refreshableInt64Constructor       = jen.Qual(refreshablePath, "NewInt64")
	refreshableInt64Ptr               = jen.Qual(refreshablePath, "Int64Ptr")
	refreshableInt64PtrConstructor    = jen.Qual(refreshablePath, "NewInt64Ptr")
	refreshableFloat64                = jen.Qual(refreshablePath, "Float64")
	refreshableFloat64Constructor     = jen.Qual(refreshablePath, "NewFloat64")
	refreshableFloat64Ptr             = jen.Qual(refreshablePath, "Float64Ptr")
	refreshableFloat64PtrConstructor  = jen.Qual(refreshablePath, "NewFloat64Ptr")
	refreshableDuration               = jen.Qual(refreshablePath, "Duration")
	refreshableDurationConstructor    = jen.Qual(refreshablePath, "NewDuration")
	refreshableDurationPtr            = jen.Qual(refreshablePath, "DurationPtr")
	refreshableDurationPtrConstructor = jen.Qual(refreshablePath, "NewDurationPtr")
)

func longestCommonPkgPathSuffix(pkgA []string, pkgB []string) int {
	if len(pkgA) == 0 || len(pkgB) == 0 {
		return 0
	}
	longestPossibleDiff := int(math.Min(float64(len(pkgA)), float64(len(pkgB))))
	for i := 0; i < longestPossibleDiff; i++ {
		if pkgA[len(pkgA)-1-i] != pkgB[len(pkgB)-1-i] {
			return i
		}
	}
	return longestPossibleDiff
}

func sanitizePackageAlias(alias string) string {
	alias = strings.ReplaceAll(alias, ".", "")
	alias = strings.ReplaceAll(alias, "-", "")
	return alias
}

// refreshableLibraryImpl returns the library implementations for the type and constructor if they exist.
// Locally defined types will result in nil return values.
func refreshableLibraryImpl(rt RefreshableType) (jenType, jenConstructor *jen.Statement) {
	switch t := rt.Type.(type) {
	case *types.Slice:
		switch elem := t.Elem().(type) {
		case *types.Basic:
			if elem.Kind() == types.String {
				return refreshableStringSlice, refreshableStringSliceConstructor
			}
		}
	case *types.Pointer:
		switch elem := t.Elem().(type) {
		case *types.Basic:
			switch elem.Kind() {
			case types.Bool:
				return refreshableBoolPtr, refreshableBoolPtrConstructor
			case types.Int:
				return refreshableIntPtr, refreshableIntPtrConstructor
			case types.Int64:
				return refreshableInt64Ptr, refreshableInt64PtrConstructor
			case types.Float64:
				return refreshableFloat64Ptr, refreshableFloat64PtrConstructor
			case types.String:
				return refreshableStringPtr, refreshableStringPtrConstructor
			}
		case *types.Named:
			if isDuration(elem) {
				return refreshableDurationPtr, refreshableDurationPtrConstructor
			}
		}
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return refreshableBool, refreshableBoolConstructor
		case types.Int:
			return refreshableInt, refreshableIntConstructor
		case types.Int64:
			return refreshableInt64, refreshableInt64Constructor
		case types.Float64:
			return refreshableFloat64, refreshableFloat64Constructor
		case types.String:
			return refreshableString, refreshableStringConstructor
		}
	case *types.Named:
		if isDuration(t) {
			return refreshableDuration, refreshableDurationConstructor
		}
	}
	return nil, nil
}

func isDuration(t *types.Named) bool {
	return t.Obj().Pkg().Path() == "time" && t.Obj().Name() == "Duration"
}
