package generic

import (
	"github.com/anubhav06/sanrakshya-cli/syft/artifact"
	"github.com/anubhav06/sanrakshya-cli/syft/file"
	"github.com/anubhav06/sanrakshya-cli/syft/linux"
	"github.com/anubhav06/sanrakshya-cli/syft/pkg"
)

type Environment struct {
	LinuxRelease *linux.Release
}

type Parser func(file.Resolver, *Environment, file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error)
