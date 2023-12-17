package haskell

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/anubhav06/sanrakshya-cli/internal/log"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/artifact"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

var _ generic.Parser = parseStackYaml

type stackYaml struct {
	ExtraDeps []string `yaml:"extra-deps"`
}

// parseStackYaml is a parser function for stack.yaml contents, returning all packages discovered.
func parseStackYaml(_ file.Resolver, _ *generic.Environment, reader file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load stack.yaml file: %w", err)
	}

	var stackFile stackYaml

	if err := yaml.Unmarshal(bytes, &stackFile); err != nil {
		log.WithFields("error", err).Tracef("failed to parse stack.yaml file %q", reader.RealPath)
		return nil, nil, nil
	}

	var pkgs []pkg.Package
	for _, dep := range stackFile.ExtraDeps {
		pkgName, pkgVersion, pkgHash := parseStackPackageEncoding(dep)
		pkgs = append(
			pkgs,
			newPackage(
				pkgName,
				pkgVersion,
				pkg.HackageStackYamlEntry{
					PkgHash: pkgHash,
				},
				reader.Location,
			),
		)
	}

	return pkgs, nil, nil
}
