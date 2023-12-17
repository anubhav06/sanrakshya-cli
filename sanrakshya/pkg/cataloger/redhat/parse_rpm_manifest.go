package redhat

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/anubhav06/sanrakshya-cli/internal/log"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/artifact"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

// Parses an RPM manifest file, as used in Mariner distroless containers, and returns the Packages listed
func parseRpmManifest(_ file.Resolver, _ *generic.Environment, reader file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	r := bufio.NewReader(reader)
	allPkgs := make([]pkg.Package, 0)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, nil, err
		}

		if line == "" {
			continue
		}

		metadata, err := newMetadataFromManifestLine(strings.TrimSuffix(line, "\n"))
		if err != nil {
			log.Warnf("unable to parse RPM manifest entry: %+v", err)
			continue
		}

		if metadata == nil {
			log.Warn("unable to parse RPM manifest entry: no metadata found")
			continue
		}

		p := newDBPackage(reader.Location, *metadata, nil, nil)

		if !pkg.IsValid(&p) {
			continue
		}

		p.SetID()
		allPkgs = append(allPkgs, p)
	}

	return allPkgs, nil, nil
}
