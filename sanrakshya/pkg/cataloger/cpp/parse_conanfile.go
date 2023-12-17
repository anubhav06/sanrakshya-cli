package cpp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/anubhav06/sanrakshya-cli/sanrakshya/artifact"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

var _ generic.Parser = parseConanfile

type Conanfile struct {
	Requires []string `toml:"requires"`
}

// parseConanfile is a parser function for conanfile.txt contents, returning all packages discovered.
func parseConanfile(_ file.Resolver, _ *generic.Environment, reader file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	r := bufio.NewReader(reader)
	inRequirements := false
	var pkgs []pkg.Package
	for {
		line, err := r.ReadString('\n')
		switch {
		case errors.Is(io.EOF, err):
			return pkgs, nil, nil
		case err != nil:
			return nil, nil, fmt.Errorf("failed to parse conanfile.txt file: %w", err)
		}

		switch {
		case strings.Contains(line, "[requires]"):
			inRequirements = true
		case strings.ContainsAny(line, "[]") || strings.HasPrefix(strings.TrimSpace(line), "#"):
			inRequirements = false
		}

		m := pkg.ConanfileEntry{
			Ref: strings.Trim(line, "\n"),
		}

		if !inRequirements {
			continue
		}

		p := newConanfilePackage(
			m,
			reader.Location.WithAnnotation(pkg.EvidenceAnnotationKey, pkg.PrimaryEvidenceAnnotation),
		)
		if p == nil {
			continue
		}

		pkgs = append(pkgs, *p)
	}
}
