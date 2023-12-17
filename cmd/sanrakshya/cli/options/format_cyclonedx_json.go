package options

import (
	"github.com/hashicorp/go-multierror"

	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/cyclonedxjson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
)

type FormatCyclonedxJSON struct {
	Pretty *bool `yaml:"pretty" json:"pretty" mapstructure:"pretty"`
}

func DefaultFormatCyclonedxJSON() FormatCyclonedxJSON {
	return FormatCyclonedxJSON{}
}

func (o FormatCyclonedxJSON) formatEncoders() ([]sbom.FormatEncoder, error) {
	var (
		encs []sbom.FormatEncoder
		errs error
	)
	for _, v := range cyclonedxjson.SupportedVersions() {
		enc, err := cyclonedxjson.NewFormatEncoderWithConfig(o.buildConfig(v))
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			encs = append(encs, enc)
		}
	}
	return encs, errs
}

func (o FormatCyclonedxJSON) buildConfig(version string) cyclonedxjson.EncoderConfig {
	var pretty bool
	if o.Pretty != nil {
		pretty = *o.Pretty
	}
	return cyclonedxjson.EncoderConfig{
		Version: version,
		Pretty:  pretty,
	}
}
