// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gotypes

import (
	"go/types"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

// FindType searches through all definitions in the package for the named type specified by typeName.
func FindType(pkg *packages.Package, typeName string) (types.Object, error) {
	for _, object := range pkg.TypesInfo.Defs {
		if object == nil {
			continue
		}
		name, ok := object.(*types.TypeName)
		if !ok {
			continue
		}
		if name.Name() == typeName {
			return name, nil
		}
	}
	return nil, errors.Errorf("package %s: did not find type %s", pkg.ID, typeName)
}
