package test2

//go:refreshable OtherStruct SuperStruct

import (
	"github.com/palantir/godel-refreshables-plugin/integration_test/testcode/test1/librarypkg"
)


type OtherStruct struct {
	FieldA string
	FieldB librarypkg.LibraryStruct
	fieldC string
}
