package sanrakshyajson

import (
	"encoding/json"
	"io"

	"github.com/anubhav06/sanrakshya-cli/internal"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
)

var _ sbom.FormatEncoder = (*encoder)(nil)

const ID sbom.FormatID = "sanrakshya-json"

type EncoderConfig struct {
	Legacy bool // transform the output to the legacy sanrakshya-json format (pre v1.0 changes, enumerated in the README.md)
	Pretty bool // don't include spaces and newlines; same as jq -c
}

type encoder struct {
	cfg EncoderConfig
}

func NewFormatEncoder() sbom.FormatEncoder {
	enc, err := NewFormatEncoderWithConfig(DefaultEncoderConfig())
	if err != nil {
		panic(err)
	}
	return enc
}

func NewFormatEncoderWithConfig(cfg EncoderConfig) (sbom.FormatEncoder, error) {
	return encoder{
		cfg: cfg,
	}, nil
}

func DefaultEncoderConfig() EncoderConfig {
	return EncoderConfig{
		Legacy: false,
		Pretty: false,
	}
}

func (e encoder) ID() sbom.FormatID {
	return ID
}

func (e encoder) Aliases() []string {
	return []string{
		"json",
		"sanrakshya",
	}
}

func (e encoder) Version() string {
	return internal.JSONSchemaVersion
}

func (e encoder) Encode(writer io.Writer, s sbom.SBOM) error {
	doc := ToFormatModel(s, e.cfg)

	enc := json.NewEncoder(writer)

	enc.SetEscapeHTML(false)

	if e.cfg.Pretty {
		enc.SetIndent("", " ")
	}

	return enc.Encode(&doc)
}
