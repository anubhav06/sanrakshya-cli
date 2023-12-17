package githubactions

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/anubhav06/sanrakshya-cli/sanrakshya/artifact"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

var _ generic.Parser = parseCompositeActionForActionUsage

type compositeActionDef struct {
	Runs compositeActionRunsDef `yaml:"runs"`
}

type compositeActionRunsDef struct {
	Steps []stepDef `yaml:"steps"`
}

func parseCompositeActionForActionUsage(_ file.Resolver, _ *generic.Environment, reader file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	contents, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read yaml composite action file: %w", err)
	}

	var ca compositeActionDef
	if err = yaml.Unmarshal(contents, &ca); err != nil {
		return nil, nil, fmt.Errorf("unable to parse yaml composite action file: %w", err)
	}

	// we use a collection to help with deduplication before raising to higher level processing
	pkgs := pkg.NewCollection()

	for _, step := range ca.Runs.Steps {
		if step.Uses == "" {
			continue
		}

		p := newPackageFromUsageStatement(step.Uses, reader.Location)
		if p != nil {
			pkgs.Add(*p)
		}
	}

	return pkgs.Sorted(), nil, nil
}
