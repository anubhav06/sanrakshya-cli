package options

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/scylladb/go-set/strset"

	"github.com/anchore/clio"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/cyclonedxjson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/sanrakshyajson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/spdxjson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/table"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
)

var _ interface {
	clio.FlagAdder
	clio.PostLoader
} = (*Output)(nil)

// Output has the standard output options sanrakshya accepts: multiple -o, --file, --template
type Output struct {
	AllowableOptions     []string `yaml:"-" json:"-" mapstructure:"-"`
	AllowMultipleOutputs bool     `yaml:"-" json:"-" mapstructure:"-"`
	AllowToFile          bool     `yaml:"-" json:"-" mapstructure:"-"`
	Outputs              []string `yaml:"output" json:"output" mapstructure:"output"` // -o, the format to use for output
	OutputFile           `yaml:",inline" json:"" mapstructure:",squash"`
	Format               `yaml:"format" json:"format" mapstructure:"format"`
}

func DefaultOutput() Output {
	return Output{
		AllowMultipleOutputs: true,
		AllowToFile:          true,
		Outputs:              []string{string(table.ID)},
		OutputFile: OutputFile{
			Enabled: true,
		},
		Format: DefaultFormat(),
	}
}

func (o *Output) PostLoad() error {
	var errs error
	for _, loader := range []clio.PostLoader{&o.OutputFile, &o.Format} {
		if err := loader.PostLoad(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (o *Output) AddFlags(flags clio.FlagSet) {
	var names []string
	for _, id := range supportedIDs() {
		names = append(names, id.String())
	}
	sort.Strings(names)

	flags.StringArrayVarP(&o.Outputs, "output", "o",
		fmt.Sprintf("report output format, formats=%v", names))
}

func (o Output) SBOMWriter() (sbom.Writer, error) {
	names := o.OutputNameSet()

	if len(o.Outputs) > 1 && !o.AllowMultipleOutputs {
		return nil, fmt.Errorf("only one output format is allowed (given %d: %s)", len(o.Outputs), names)
	}

	encoders, err := o.Encoders()
	if err != nil {
		return nil, err
	}

	if !o.AllowToFile {
		for _, opt := range o.Outputs {
			if strings.Contains(opt, "=") {
				return nil, fmt.Errorf("file output is not allowed ('-o format=path' should be '-o format')")
			}
		}
	}

	return makeSBOMWriter(o.Outputs, o.File, encoders)
}

func (o Output) OutputNameSet() *strset.Set {
	names := strset.New()
	for _, output := range o.Outputs {
		fields := strings.Split(output, "=")
		names.Add(fields[0])
	}

	return names
}

func supportedIDs() []sbom.FormatID {
	encs := []sbom.FormatID{
		// encoders that support a single version
		sanrakshyajson.ID,
		table.ID,

		// encoders that support multiple versions
		cyclonedxjson.ID,
		spdxjson.ID,
	}

	return encs
}
