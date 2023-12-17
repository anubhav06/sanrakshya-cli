package github

import (
	"encoding/json"
	"io"

	"github.com/anubhav06/sanrakshya-cli/sanrakshya/format/github/internal/model"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
)

const ID sbom.FormatID = "github-json"

type encoder struct {
}

func NewFormatEncoder() sbom.FormatEncoder {
	return encoder{}
}

func (e encoder) ID() sbom.FormatID {
	return ID
}

func (e encoder) Aliases() []string {
	return []string{
		"github",
	}
}

func (e encoder) Version() string {
	return sbom.AnyVersion
}

func (e encoder) Encode(writer io.Writer, s sbom.SBOM) error {
	bom := model.ToGithubModel(&s)

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	return enc.Encode(bom)
}
