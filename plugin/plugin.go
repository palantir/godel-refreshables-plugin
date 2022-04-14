// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/palantir/godel-refreshables-plugin/config"
	"github.com/palantir/godel-refreshables-plugin/plugin/generator"
	"github.com/palantir/godel-refreshables-plugin/plugin/gotypes"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

const (
	fullLoadMode = packages.NeedName | packages.NeedDeps | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedFiles
)

var fset = token.NewFileSet()

func Run(projectDir string, cfg config.Config, verify bool) error {
	for pkgPath, pkgCfg := range cfg.Refreshables {
		if err := renderRefreshableTypesFile(projectDir, cfg.ImportAliases, pkgPath, pkgCfg.Types, pkgCfg.Output, verify); err != nil {
			return errors.Wrap(err, pkgPath)
		}
	}
	return nil
}

func renderRefreshableTypesFile(projectDir string, importAliases map[string]string, pkgPath string, typeNames []string, outputFile string, verify bool) error {
	pkg, err := loadPackage(projectDir, pkgPath)
	if err != nil {
		return err
	}

	outputFile, outputPackagePath, outputPackageName, err := getOutputSpec(projectDir, pkg, outputFile)
	if err != nil {
		return err
	}

	// Collect all nested types -> load packages for all nested types -> resolve to remote or local refreshable -> generate code

	typeObjects := make([]types.Type, len(typeNames))
	for i, typeName := range typeNames {
		typeObj, err := gotypes.FindType(pkg, typeName)
		if err != nil {
			return err
		}
		typeObjects[i] = typeObj.Type()
	}
	typeObjects, err = gotypes.FlattenTypes(typeObjects...)
	if err != nil {
		return err
	}

	refreshableTypes, err := generator.NewRefreshableTypes(pkg, typeObjects)
	if err != nil {
		return err
	}

	file, err := generator.GenerateRefreshableFile(outputPackagePath, outputPackageName, refreshableTypes)
	if err != nil {
		return err
	}
	for path, alias := range importAliases {
		file.ImportAlias(path, alias)
	}
	buf := &bytes.Buffer{}
	if err := file.Render(buf); err != nil {
		return err
	}
	outputBytes, err := imports.Process(outputFile, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	if verify {
		existing, err := ioutil.ReadFile(outputFile)
		if os.IsNotExist(err) {
			return errors.Wrap(err, "regenerate refreshables output")
		}
		if !bytes.Equal(existing, outputBytes) {
			return errors.Errorf("regenerate refreshables output: outdated file %s", outputFile)
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return errors.Wrap(err, "create outputFile parent directories")
		}
		if err := ioutil.WriteFile(outputFile, outputBytes, 0644); err != nil {
			return errors.Wrap(err, "write outputFile")
		}
	}
	return nil
}

func loadPackage(projectDir string, pkgPath string) (*packages.Package, error) {
	pkg, err := loadSinglePackage(projectDir, pkgPath, fullLoadMode)
	if err != nil {
		return nil, err
	}
	if err := validatePackage(pkg); err != nil {
		return nil, err
	}
	return pkg, nil
}

func loadSinglePackage(projectDir string, pkgPath string, mode packages.LoadMode) (*packages.Package, error) {
	loadedPackages, err := packages.Load(&packages.Config{
		Mode: mode,
		Dir:  projectDir,
		Fset: fset,
	}, pkgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: failed to load package", pkgPath)
	}
	if len(loadedPackages) != 1 {
		return nil, errors.Errorf("%s: expected exactly one loaded package, got %d", pkgPath, len(loadedPackages))
	}
	pkg := loadedPackages[0]
	return pkg, nil
}

func validatePackage(pkg *packages.Package) error {
	if pkg == nil {
		return errors.Errorf("nil package")
	}
	if len(pkg.Errors) > 0 {
		errs := ""
		for _, e := range pkg.Errors {
			errs += "\n" + e.Error()
		}
		return errors.Errorf("failed to load package %s:%s", pkg.PkgPath, errs)
	}
	if pkg.IllTyped {
		return errors.Errorf("package %s was ill-typed", pkg.PkgPath)
	}
	return nil
}

// getOutputSpec determines where the pkg's generated refreshables file will be written and its go package metadata.
// If outputFile is empty, the default location within pkg will be used.
// If outputFile is specified, it must be a go file within projectDir.
func getOutputSpec(projectDir string, pkg *packages.Package, outputFile string) (outputFilename, outputPkgPath, outputPkgName string, err error) {
	if outputFile == "" {
		if pkg.Module != nil && pkg.Module.Dir != projectDir {
			return "", "", "", errors.Errorf("output destination required for packages outside local module")
		}
		// this is a local package, generate into the package directory
		if len(pkg.GoFiles) == 0 {
			return "", "", "", errors.Errorf("pkg %s has no go files", pkg.PkgPath)
		}
		file := filepath.Join(filepath.Dir(pkg.GoFiles[0]), "zz_generated_refreshables.go")
		return file, pkg.PkgPath, pkg.Name, nil
	}

	if filepath.Ext(outputFile) != ".go" {
		return "", "", "", errors.Errorf("Output %q file extension must be .go", outputFile)
	}
	if filepath.IsAbs(outputFile) {
		return "", "", "", errors.Errorf("Output %q must be a relative path", outputFile)
	}
	if strings.HasPrefix(outputFile, ".."+string(filepath.Separator)) {
		return "", "", "", errors.Errorf("Output %q must exist within project directory %q", outputFile, projectDir)
	}
	outputFilename = filepath.Join(projectDir, outputFile)
	outputPkgPath = "./" + filepath.Dir(outputFile) // Add ./ so go doesn't treat it as a normal package path
	outputPkgName = filepath.Base(filepath.Dir(outputFilename))

	// Try to load outputDir as a package; if it fails, we'll use the derived values instead.
	if outputPkg, err := loadSinglePackage(projectDir, outputPkgPath, packages.NeedName); err == nil {
		outputPkgPath = outputPkg.PkgPath
		if outputPkg.Name != "" {
			outputPkgName = outputPkg.Name
		}
	}

	return outputFilename, outputPkgPath, outputPkgName, nil
}
