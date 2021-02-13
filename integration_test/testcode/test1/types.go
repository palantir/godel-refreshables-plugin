package test1

//go:refreshable OtherStruct SuperStruct

import (
	"time"
	"github.com/palantir/godel-refreshables-plugin/integration_test/testcode/test1/librarypkg"
)


type OtherStruct struct {
	FieldA string
	FieldB InnerStruct
}

type InnerStruct struct {
	InnerFieldA int
	InnerFieldB time.Duration
}

type SuperStruct struct {
	// Primitives and Variants
	String                    string
	OptionalString            *string
	SliceString               []string
	ArrayString               [8]string
	StringString              map[string]string
	StringAlias               StringAlias
	OptionalStringAlias       OptionalStringAlias
	DoubleOptionalStringAlias *OptionalStringAlias

	Int                    int
	OptionalInt            *int
	SliceInt               []int
	ArrayInt               [8]int
	IntInt                 map[int]int
	IntAlias               IntAlias
	OptionalIntAlias       OptionalIntAlias
	DoubleOptionalIntAlias *OptionalIntAlias

	Duration                    time.Duration
	OptionalDuration            *time.Duration
	SliceDuration               []time.Duration
	ArrayDuration               [8]time.Duration
	DurationDuration            map[time.Duration]time.Duration
	DurationAlias               DurationAlias
	OptionalDurationAlias       OptionalDurationAlias
	DoubleOptionalDurationAlias *OptionalDurationAlias

	// Local types
	NestedStruct
	NamedNestedStruct         NestedStruct
	OptionalNestedStruct      *NestedStruct
	SliceNestedStruct         []NestedStruct
	ArrayNestedStruct         [8]NestedStruct
	NestedStructNestedStruct  map[NestedStruct]NestedStruct
	NestedStructAlias         NestedStructAlias
	OptionalNestedStructAlias OptionalNestedStructAlias
	//TODO: this doesn't work due to need to dereference before accessing struct fields
	// DoubleOptionalNestedStructAlias *OptionalNestedStructAlias

	// Imported Types
	librarypkg.LibraryStruct
	NamedLibraryStruct         librarypkg.LibraryStruct
	OptionalLibraryStruct      *librarypkg.LibraryStruct
	SliceLibraryStruct         []librarypkg.LibraryStruct
	ArrayLibraryStruct         [8]librarypkg.LibraryStruct
	LibraryStructLibraryStruct map[librarypkg.LibraryStruct]librarypkg.LibraryStruct
	LibraryStructAlias         LibraryStructAlias
	OptionalLibraryStructAlias OptionalLibraryStructAlias
}

type StringAlias string

type OptionalStringAlias *StringAlias

type IntAlias int

type OptionalIntAlias *IntAlias

type DurationAlias time.Duration

type OptionalDurationAlias *DurationAlias

type NestedStruct struct {
	FieldA string
	FieldB int
}

type NestedStructBC struct {
	FieldB string
	FieldC int
}

type NestedStructAlias NestedStruct

type OptionalNestedStructAlias *NestedStructAlias

type LibraryStructAlias librarypkg.LibraryStruct

type OptionalLibraryStructAlias *LibraryStructAlias
