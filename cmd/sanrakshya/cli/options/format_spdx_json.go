package options

import (
	"github.com/hashicorp/go-multierror"

	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/spdxjson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
)

type FormatSPDXJSON struct {
	Pretty *bool `yaml:"pretty" json:"pretty" mapstructure:"pretty"`
}

func DefaultFormatSPDXJSON() FormatSPDXJSON {
	return FormatSPDXJSON{}
}

func (o FormatSPDXJSON) formatEncoders() ([]sbom.FormatEncoder, error) {
	var (
		encs []sbom.FormatEncoder
		errs error
	)
	for _, v := range spdxjson.SupportedVersions() {
		enc, err := spdxjson.NewFormatEncoderWithConfig(o.buildConfig(v))
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			encs = append(encs, enc)
		}
	}
	return encs, errs
}

func (o FormatSPDXJSON) buildConfig(v string) spdxjson.EncoderConfig {
	var pretty bool
	if o.Pretty != nil {
		pretty = *o.Pretty
	}
	return spdxjson.EncoderConfig{
		Version: v,
		Pretty:  pretty,
	}
}
