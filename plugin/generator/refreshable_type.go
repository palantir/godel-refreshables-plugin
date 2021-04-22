// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package generator

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/goastwriter/decl"
	"golang.org/x/tools/go/packages"
)

const (
	consumerFuncName    = "consumer"
	mapFuncName         = "mapFn"
	unsubscribeFuncName = "unsubscribe"
)

var (
	unsubscribeFunc = jen.Params(jen.Id(unsubscribeFuncName).Func().Params())
)

func GenerateRefreshableFile(targetPackagePath, targetPackageName string, refreshableTypes RefreshableTypes) *jen.File {
	f := jen.NewFilePathName(targetPackagePath, targetPackageName)
	f.HeaderComment("Generated by godel-refreshable-plugin: do not edit.")

	for _, rt := range refreshableTypes {
		if libraryType, _ := refreshableLibraryImpl(rt); libraryType != nil {
			// Don't generate code for the types that have library implementations
			continue
		}
		f.Type().Id(rt.ifaceTypeString()).Interface(rt.getJenInterfaceMethods(refreshableTypes)...).Line()
		f.Type().Id(rt.implTypeString()).Struct(refreshable)
		f.Add(rt.getJenImplementation(refreshableTypes)...)
		f.Line()
	}
	return f
}

// RefreshableTypes is a container for RefreshableType which provides some convenience functions
type RefreshableTypes []RefreshableType

func NewRefreshableTypes(targetPackage *packages.Package, typeSet []types.Type) (RefreshableTypes, error) {
	var refreshableTypes RefreshableTypes
	for _, typ := range typeSet {
		rt := newRefreshableType(targetPackage, typ, typeSet)
		refreshableTypes = append(refreshableTypes, rt)
	}
	return refreshableTypes, nil
}

// Imports returns an unordered list of unique imports for the analyzed fields
func (t RefreshableTypes) Imports() decl.Imports {
	var imports decl.Imports

	importSet := make(map[decl.Import]struct{})
	for _, analyzedField := range t {
		for _, i := range analyzedField.Imports {
			importSet[*i] = struct{}{}
		}
	}

	for i := range importSet {
		iCopy := i
		imports = append(imports, &iCopy)
	}

	return imports
}

func (t RefreshableTypes) forType(typ types.Type) RefreshableType {
	for _, rt := range t {
		if types.Identical(typ, rt.Type) {
			return rt
		}
	}
	return RefreshableType{Type: typ}
}

// A RefreshableType contains all the necessary information to generate an interface and implementation for the
// contained internal type. It is expected that a constructed refreshable type has already handled any potential naming
// collisions, so users of this type are safe to use the expressions and declarations returned from it's functions if
// all the RefreshableTypes used they same refreshableTypeGenerator.
type RefreshableType struct {
	Type         types.Type
	Imports      map[types.Type]*decl.Import
	OverrideName string
}

func newRefreshableType(targetPackage *packages.Package, typ types.Type, typeSet []types.Type) RefreshableType {
	namedType, ok := typ.(*types.Named)
	if !ok {
		// Nothing to analyze for primitives
		return RefreshableType{
			Type: typ,
		}
	}
	typePkg := namedType.Obj().Pkg()

	// local package means we can use simple name
	if typePkg.Path() == targetPackage.PkgPath {
		return RefreshableType{
			Type: typ,
		}
	}

	// If we are this far, we have a type external to the package and need to see if it must be renamed

	pathSlice := strings.Split(typePkg.Path(), "/")
	pkgPathSuffix := 0

	rt := RefreshableType{
		Type: typ,
		Imports: map[types.Type]*decl.Import{
			typ: {Path: typePkg.Path(), Alias: typePkg.Name()},
		},
	}

	for _, otherTyp := range typeSet {
		otherNamedType, ok := otherTyp.(*types.Named)
		if !ok {
			continue
		}
		otherTypePkg := otherNamedType.Obj().Pkg()
		// Skip comparing fields against other fields from the same package... the compiler already prevents this
		// naming collision for us
		if otherTypePkg == typePkg {
			continue
		}
		otherPathSlice := strings.Split(otherTypePkg.Path(), "/")

		// Determine if there is a package collision, and how far back in the package chain we have to go in order
		// to get a unique path for it.
		// i.e. (github.com/asdf/foo/bar/baz/quux, github.com/qwer/foo/bar/baz/quux) -> 4
		tempPkgPathSuffix := longestCommonPkgPathSuffix(pathSlice, otherPathSlice)
		if tempPkgPathSuffix > pkgPathSuffix {
			pkgPathSuffix = tempPkgPathSuffix
		}

		// Determine if there is a struct naming collision, in which case we update the field names to include their
		// aliased package names, preventing
		//
		// package github.com/asdf/foo.Config -> AsdfFooConfig
		// package github.com/qwer/foo.Config -> QwerFooConfig
		// TODO: handle ptr/slice/map names here??
		if nameFromType(typ) == nameFromType(otherTyp) {
			b := strings.Builder{}
			for _, str := range pathSlice[len(pathSlice)-1-pkgPathSuffix:] {
				b.WriteString(sanitizePackageAlias(strings.Title(str)))
			}
			b.WriteString(nameFromType(typ))
			rt.OverrideName = b.String()
		}
	}

	if pkgPathSuffix != 0 {
		rt.Imports[typ].Alias = sanitizePackageAlias(strings.Join(pathSlice[len(pathSlice)-1-pkgPathSuffix:], ""))
	}
	return rt
}

func (rt RefreshableType) uniqueName() string {
	if rt.OverrideName != "" {
		return rt.OverrideName
	}
	return strings.Title(nameFromType(rt.Type))
}

func (rt RefreshableType) getJenInterfaceMethods(typeSet RefreshableTypes) []jen.Code {
	methods := []jen.Code{
		refreshable,
		rt.currentFunc(),
		rt.mapFunc(),
		rt.subscribeFunc(),
		jen.Line(),
	}
	methods = append(methods, interfaceMethods(rt.Type, typeSet)...)

	return methods
}

func (rt RefreshableType) currentFunc() jen.Code {
	return jen.Id(rt.currentFuncString()).Params().Add(jenType(rt.Type))
}

func (rt RefreshableType) mapFunc() jen.Code {
	return jen.Id(rt.mapFuncString()).Params(jen.Func().Params(jenType(rt.Type)).Add(jen.Interface())).Add(refreshable)
}

func (rt RefreshableType) subscribeFunc() jen.Code {
	return jen.Id(rt.subscribeFuncString()).Params(jen.Func().Params(jenType(rt.Type))).Add(jen.Params(jen.Id("unsubscribe").Func().Params()))
}

func interfaceMethods(typ types.Type, typeSet RefreshableTypes) []jen.Code {
	switch underlying := typ.Underlying().(type) {
	case *types.Pointer:
		return interfaceMethods(underlying.Elem(), typeSet)
	case *types.Struct:
		var methods []jen.Code
		for i := 0; i < underlying.NumFields(); i++ {
			field := underlying.Field(i)
			if field.Exported() {
				methods = append(methods, jen.Id(field.Name()).Params().Add(refreshableJenType(typeSet.forType(field.Type()))))
			}
		}
		return methods
	}
	return nil
}

func refreshableJenType(rt RefreshableType) jen.Code {
	libraryType, _ := refreshableLibraryImpl(rt)
	if libraryType != nil {
		return libraryType
	}
	return jen.Id(rt.ifaceTypeString())
}

func refreshableJenTypeConstructor(rt RefreshableType) jen.Code {
	_, libraryConstructor := refreshableLibraryImpl(rt)
	if libraryConstructor != nil {
		return libraryConstructor
	}
	return jen.Id(fmt.Sprintf("New%s", rt.implTypeString()))
}

func jenType(typ types.Type) jen.Code {
	switch t := typ.(type) {
	case *types.Pointer:
		return jen.Op("*").Add(jenName(typ))
	case *types.Slice:
		return jen.Index().Add(jenName(typ))
	case *types.Map:
		return jen.Map(jenName(t.Key())).Add(jenName(t.Elem()))
	default:
		return jenName(typ)
	}
}

func jenName(typ types.Type) jen.Code {
	switch t := typ.(type) {
	case *types.Interface:
		return jen.Interface()
	case *types.Map:
		return jen.Map(jen.Add(jenName(t.Key()))).Add(jenName(t.Elem()))
	case *types.Pointer:
		return jenName(t.Elem())
	case *types.Array:
		return jenName(t.Elem())
	case *types.Slice:
		return jenName(t.Elem())
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.UntypedBool:
			return jen.Bool()
		case types.Int, types.UntypedInt:
			return jen.Int()
		case types.Int8:
			return jen.Int8()
		case types.Int16:
			return jen.Int16()
		case types.Int32:
			return jen.Int32()
		case types.Int64:
			if typ.String() == "time.Duration" {
				return jen.Qual("time", "Duration")
			}
			return jen.Int64()
		case types.Uint:
			return jen.Uint()
		case types.Uint8:
			return jen.Uint8()
		case types.Uint16:
			return jen.Uint16()
		case types.Uint32:
			return jen.Uint32()
		case types.Uint64:
			return jen.Uint64()
		case types.Uintptr:
			return jen.Uintptr()
		case types.Float32:
			return jen.Float32()
		case types.Float64, types.UntypedFloat:
			return jen.Float64()
		case types.Complex64:
			return jen.Complex64()
		case types.Complex128, types.UntypedComplex:
			return jen.Complex128()
		case types.String, types.UntypedString:
			return jen.String()
		case types.UntypedRune:
			return jen.Rune()
		case types.UntypedNil:
			return jen.Nil()
		}
	case *types.Named:
		return jen.Qual(t.Obj().Pkg().Path(), t.Obj().Name())
	}
	var name string
	qualifiedName := strings.Split(nameFromType(typ), ".")
	if len(name) > 1 {
		name = qualifiedName[1]
	} else {
		name = qualifiedName[0]
	}
	return jen.Id(name)
}

func nameFromType(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Basic:
		return t.Name()
	case *types.Interface:
		return "Any"
	case *types.Named:
		return strings.Title(t.Obj().Name())
	case *types.Map:
		k := strings.Title(nameFromType(t.Key()))
		v := strings.Title(nameFromType(t.Elem()))
		return fmt.Sprintf("%sTo%s", k, v)
	case *types.Pointer:
		return fmt.Sprintf("%sPtr", strings.Title(nameFromType(t.Elem())))
	case *types.Slice:
		return fmt.Sprintf("%sSlice", strings.Title(nameFromType(t.Elem())))
	case *types.Array:
		return fmt.Sprintf("%sArray", strings.Title(nameFromType(t.Elem())))
	default:
		return strings.Title(t.String())
	}
}

func (rt RefreshableType) getJenImplementation(typeSet RefreshableTypes) []jen.Code {
	var methods []jen.Code
	methods = append(methods, rt.constructor(), jen.Line(), jen.Line())
	methods = append(methods, rt.typedCurrentImpl(), jen.Line(), jen.Line())
	methods = append(methods, rt.typedMapImpl(), jen.Line(), jen.Line())
	methods = append(methods, rt.typedSubscribeImpl(), jen.Line(), jen.Line())
	methods = append(methods, rt.implementationMethods(rt.Type, typeSet)...)

	return methods
}

func (rt RefreshableType) constructor() jen.Code {
	return jen.Func().Id(fmt.Sprintf("New%s", rt.implTypeString())).Params(jen.Id("in").Add(refreshable)).Add(jen.Id(rt.implTypeString())).Block(
		jen.Return(jen.Id(rt.implTypeString()).Values(jen.Dict{jen.Id("Refreshable"): jen.Id("in")})),
	)
}

func (rt RefreshableType) typedCurrentImpl() jen.Code {
	return jen.Func().Params(
		jen.Id("r").Id(rt.implTypeString()),
	).Id(rt.currentFuncString()).Params().Add(jenType(rt.Type)).Block(
		jen.Return(jen.Id("r").Dot("Current").Call().Assert(jenType(rt.Type))),
	)
}

func (rt RefreshableType) typedMapImpl() jen.Code {
	return jen.Func().Params(
		jen.Id("r").Id(rt.implTypeString()),
	).Id(rt.mapFuncString()).Params(rt.typedMapFn()).Add(refreshable).Block(
		jen.Return(jen.Id("r").Dot("Map").Call(
			jen.Func().Params(jen.Id("i").Interface()).Interface().Block(
				jen.Return(jen.Id(mapFuncName).Params(jen.Id("i").Assert(jenType(rt.Type)))),
			),
		)),
	)
}

func (rt RefreshableType) typedSubscribeImpl() jen.Code {
	return jen.Func().Params(
		jen.Id("r").Id(rt.implTypeString()),
	).Id(rt.subscribeFuncString()).Params(rt.typedSubscribeFn()).Add(unsubscribeFunc).Block(
		jen.Return(jen.Id("r").Dot("Subscribe").Call(
			jen.Func().Params(jen.Id("i").Interface()).Block(
				jen.Id(consumerFuncName).Params(jen.Id("i").Assert(jenType(rt.Type))),
			),
		)),
	)
}

func (rt RefreshableType) implementationMethods(t types.Type, typeSet RefreshableTypes) []jen.Code {
	switch underlying := t.Underlying().(type) {
	case *types.Struct:
		var methods []jen.Code
		for i := 0; i < underlying.NumFields(); i++ {
			field := underlying.Field(i)
			typ := typeSet.forType(field.Type())
			if field.Exported() {
				methods = append(methods, rt.jenImplementationForField(underlying.Field(i), typ), jen.Line(), jen.Line())
			}
		}
		return methods
	case *types.Pointer:
		return rt.implementationMethods(underlying.Elem(), typeSet)
	}
	return nil
}

func (rt RefreshableType) jenImplementationForField(field *types.Var, typ RefreshableType) jen.Code {
	return jen.Func().Params(
		jen.Id("r").Id(rt.implTypeString()),
	).Id(field.Name()).Params().Add(refreshableJenType(typ)).Block(
		jen.Return(refreshableJenTypeConstructor(typ)).Params(jen.Id("r").Dot(rt.mapFuncString()).Params(
			jen.Func().Params(jen.Id("i").Add(jenType(rt.Type))).Interface().Block(
				jen.Return(jen.Id("i").Dot(field.Name())),
			),
		)),
	)
}

func (rt RefreshableType) ifaceTypeString() string {
	return fmt.Sprintf("Refreshable%s", rt.uniqueName())
}

func (rt RefreshableType) implTypeString() string {
	return fmt.Sprintf("Refreshing%s", rt.uniqueName())
}

func (rt RefreshableType) currentFuncString() string {
	return fmt.Sprintf("Current%s", rt.uniqueName())
}

func (rt RefreshableType) mapFuncString() string {
	return fmt.Sprintf("Map%s", rt.uniqueName())
}

func (rt RefreshableType) subscribeFuncString() string {
	return fmt.Sprintf("SubscribeTo%s", rt.uniqueName())
}

func (rt RefreshableType) typedMapFn() jen.Code {
	return jen.Id(mapFuncName).Func().Params(jenType(rt.Type)).Add(jen.Interface())
}

func (rt RefreshableType) typedSubscribeFn() jen.Code {
	return jen.Id(consumerFuncName).Func().Params(jenType(rt.Type))
}
