// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"go/token"
	"go/types"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/palantir/godel-refreshables-plugin/config"
	"github.com/palantir/godel-refreshables-plugin/plugin/generator"
	"github.com/palantir/godel-refreshables-plugin/plugin/gotypes"
	"github.com/palantir/pkg/codegenfiles"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

const (
	fullLoadMode = packages.NeedName | packages.NeedDeps | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedFiles
)

var fset = token.NewFileSet()

func Run(projectDir string, cfg config.Config, verify bool) error {
	out := codegenfiles.NewOutput()
	// packages are rendered in a deterministic order so that the first configuration error reported does
	// not depend on map iteration order
	for _, pkgPath := range slices.Sorted(maps.Keys(cfg.Refreshables)) {
		pkgCfg := cfg.Refreshables[pkgPath]
		outputFile, outputBytes, err := renderRefreshableTypesFile(projectDir, cfg.ImportAliases, pkgPath, pkgCfg.Types, pkgCfg.Output)
		if err != nil {
			return errors.Wrap(err, pkgPath)
		}
		out.Add(outputFile, outputBytes)
	}

	// Every configured package writes into the project directory, so one project reconciles all of them.
	// That is what makes two entries resolving to the same output file an error rather than a
	// last-writer-wins race, and what keeps an output path from escaping the project directory.
	//
	// DeleteStale is deliberately off: the plugin's output locations come entirely from configuration, so
	// it has no record of where an earlier run wrote, and the only name it could search for --
	// zz_generated_refreshables.go -- is also the name under which consumers vendor this plugin's output
	// from their dependencies. Deleting on that name would remove vendored source.
	p := &codegenfiles.Project{Dir: projectDir}
	changes, err := p.Plan(out)
	if err != nil {
		return err
	}
	if verify {
		return changes.Err()
	}
	return changes.Apply()
}

// renderRefreshableTypesFile returns the absolute path of the file generated for pkgPath and its content.
func renderRefreshableTypesFile(projectDir string, importAliases map[string]string, pkgPath string, typeNames []string, outputFile string) (string, []byte, error) {
	pkg, err := loadPackage(projectDir, pkgPath)
	if err != nil {
		return "", nil, err
	}

	outputFile, outputPackagePath, outputPackageName, err := getOutputSpec(projectDir, pkg, outputFile)
	if err != nil {
		return "", nil, err
	}

	// Collect all nested types -> load packages for all nested types -> resolve to remote or local refreshable -> generate code

	typeObjects := make([]types.Type, len(typeNames))
	for i, typeName := range typeNames {
		typeObj, err := gotypes.FindType(pkg, typeName)
		if err != nil {
			return "", nil, err
		}
		typeObjects[i] = typeObj.Type()
	}
	typeObjects, err = gotypes.FlattenTypes(typeObjects...)
	if err != nil {
		return "", nil, err
	}

	refreshableTypes, err := generator.NewRefreshableTypes(pkg, typeObjects)
	if err != nil {
		return "", nil, err
	}

	file, err := generator.GenerateRefreshableFile(outputPackagePath, outputPackageName, refreshableTypes)
	if err != nil {
		return "", nil, err
	}
	for path, alias := range importAliases {
		file.ImportAlias(path, alias)
	}
	buf := &bytes.Buffer{}
	if err := file.Render(buf); err != nil {
		return "", nil, err
	}
	outputBytes, err := imports.Process(outputFile, buf.Bytes(), nil)
	if err != nil {
		return "", nil, err
	}
	return outputFile, outputBytes, nil
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
		var errs strings.Builder
		for _, e := range pkg.Errors {
			errs.WriteString("\n" + e.Error())
		}
		return errors.Errorf("failed to load package %s:%s", pkg.PkgPath, errs.String())
	}
	if pkg.IllTyped {
		return errors.Errorf("package %s was ill-typed", pkg.PkgPath)
	}
	return nil
}

// getOutputSpec determines where the pkg's generated refreshables file will be written and its go package metadata.
// If outputFile is empty, the default location within pkg will be used.
// If outputFile is specified, it must be a relative go file path; that it resolves within projectDir is
// enforced when the output is reconciled, which a check on the path text cannot do reliably.
// The returned filename is absolute.
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
	outputFilename, err = filepath.Abs(filepath.Join(projectDir, outputFile))
	if err != nil {
		return "", "", "", errors.Wrapf(err, "Output %q", outputFile)
	}
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
