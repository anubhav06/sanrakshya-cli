package options

import (
	"github.com/anubhav06/sanrakshya-cli/syft/format/syftjson"
	"github.com/anubhav06/sanrakshya-cli/syft/sbom"
)

type FormatSyftJSON struct {
	Legacy bool  `yaml:"legacy" json:"legacy" mapstructure:"legacy"`
	Pretty *bool `yaml:"pretty" json:"pretty" mapstructure:"pretty"`
}

func DefaultFormatJSON() FormatSyftJSON {
	return FormatSyftJSON{
		Legacy: false,
	}
}

func (o FormatSyftJSON) formatEncoders() ([]sbom.FormatEncoder, error) {
	enc, err := syftjson.NewFormatEncoderWithConfig(o.buildConfig())
	return []sbom.FormatEncoder{enc}, err
}

func (o FormatSyftJSON) buildConfig() syftjson.EncoderConfig {
	var pretty bool
	if o.Pretty != nil {
		pretty = *o.Pretty
	}
	return syftjson.EncoderConfig{
		Legacy: o.Legacy,
		Pretty: pretty,
	}
}
