package options

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/sanrakshyajson"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
)

type FormatsanrakshyaJSON struct {
	Legacy bool  `yaml:"legacy" json:"legacy" mapstructure:"legacy"`
	Pretty *bool `yaml:"pretty" json:"pretty" mapstructure:"pretty"`
}

func DefaultFormatJSON() FormatsanrakshyaJSON {
	return FormatsanrakshyaJSON{
		Legacy: false,
	}
}

func (o FormatsanrakshyaJSON) formatEncoders() ([]sbom.FormatEncoder, error) {
	enc, err := sanrakshyajson.NewFormatEncoderWithConfig(o.buildConfig())
	return []sbom.FormatEncoder{enc}, err
}

func (o FormatsanrakshyaJSON) buildConfig() sanrakshyajson.EncoderConfig {
	var pretty bool
	if o.Pretty != nil {
		pretty = *o.Pretty
	}
	return sanrakshyajson.EncoderConfig{
		Legacy: o.Legacy,
		Pretty: pretty,
	}
}
