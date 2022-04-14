package test3

import (
	"github.com/palantir/godel-refreshables-plugin/integration_test/testcode/test3/librarypkg"
)

type MainStruct struct {
	ExcludedString string    `yaml:"excluded-string" refreshables:",exclude"`
	IncludedString string    `yaml:"included-string" refreshables:"custom-name"`
	Sub            SubStruct `yaml:"sub"`
}

type SubStruct struct {
	IncludedInt     int                            `yaml:"included-int"`
	ExcludedLibrary librarypkg.CustomLibraryStruct `yaml:"excluded-library" refreshables:",exclude"`
}
