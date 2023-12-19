package options

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/anchore/clio"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/cyclonedxjson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/sanrakshyajson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/spdxjson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/table"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
)

var _ clio.PostLoader = (*Format)(nil)

// Format contains all user configuration for output formatting.
type Format struct {
	Pretty         *bool                `yaml:"pretty" json:"pretty" mapstructure:"pretty"`
	SanrakshyaJSON FormatsanrakshyaJSON `yaml:"json" json:"json" mapstructure:"json"`
	SPDXJSON       FormatSPDXJSON       `yaml:"spdx-json" json:"spdx-json" mapstructure:"spdx-json"`
	CyclonedxJSON  FormatCyclonedxJSON  `yaml:"cyclonedx-json" json:"cyclonedx-json" mapstructure:"cyclonedx-json"`
}

func (o *Format) PostLoad() error {
	o.SanrakshyaJSON.Pretty = multiLevelOption[bool](false, o.Pretty, o.SanrakshyaJSON.Pretty)
	o.SPDXJSON.Pretty = multiLevelOption[bool](false, o.Pretty, o.SPDXJSON.Pretty)
	o.CyclonedxJSON.Pretty = multiLevelOption[bool](false, o.Pretty, o.CyclonedxJSON.Pretty)

	return nil
}

func DefaultFormat() Format {
	return Format{
		SanrakshyaJSON: DefaultFormatJSON(),
		SPDXJSON:       DefaultFormatSPDXJSON(),
		CyclonedxJSON:  DefaultFormatCyclonedxJSON(),
	}
}

func (o *Format) Encoders() ([]sbom.FormatEncoder, error) {
	// setup all encoders based on the configuration
	var list encoderList

	// in the future there will be application configuration options that can be used to set the default output format
	list.addWithErr(sanrakshyajson.ID)(o.SanrakshyaJSON.formatEncoders())
	list.add(table.ID)(table.NewFormatEncoder())
	list.addWithErr(cyclonedxjson.ID)(o.CyclonedxJSON.formatEncoders())
	list.addWithErr(spdxjson.ID)(o.SPDXJSON.formatEncoders())

	return list.encoders, list.err
}

type encoderList struct {
	encoders []sbom.FormatEncoder
	err      error
}

func (l *encoderList) addWithErr(name sbom.FormatID) func([]sbom.FormatEncoder, error) {
	return func(encs []sbom.FormatEncoder, err error) {
		if err != nil {
			l.err = multierror.Append(l.err, fmt.Errorf("unable to configure %q format encoder: %w", name, err))
			return
		}
		for _, enc := range encs {
			if enc == nil {
				l.err = multierror.Append(l.err, fmt.Errorf("unable to configure %q format encoder: nil encoder returned", name))
				continue
			}
			l.encoders = append(l.encoders, enc)
		}
	}
}

func (l *encoderList) add(name sbom.FormatID) func(...sbom.FormatEncoder) {
	return func(encs ...sbom.FormatEncoder) {
		for _, enc := range encs {
			if enc == nil {
				l.err = multierror.Append(l.err, fmt.Errorf("unable to configure %q format encoder: nil encoder returned", name))
				continue
			}
			l.encoders = append(l.encoders, enc)
		}
	}
}

func multiLevelOption[T any](defaultValue T, option ...*T) *T {
	result := defaultValue
	for _, opt := range option {
		if opt != nil {
			result = *opt
		}
	}
	return &result
}
