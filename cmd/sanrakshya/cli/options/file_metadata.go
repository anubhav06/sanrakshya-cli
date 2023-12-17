package options

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/source"
)

type fileMetadata struct {
	Cataloger scope    `yaml:"cataloger" json:"cataloger" mapstructure:"cataloger"`
	Digests   []string `yaml:"digests" json:"digests" mapstructure:"digests"`
}

func defaultFileMetadata() fileMetadata {
	return fileMetadata{
		Cataloger: scope{
			Scope: source.SquashedScope.String(),
		},
		Digests: []string{"sha256"},
	}
}
