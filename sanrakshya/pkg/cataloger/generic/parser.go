package generic

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/artifact"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/linux"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
)

type Environment struct {
	LinuxRelease *linux.Release
}

type Parser func(file.Resolver, *Environment, file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error)
