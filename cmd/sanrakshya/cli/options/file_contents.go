package options

import (
	"github.com/anubhav06/sanrakshya-cli/internal/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/source"
)

type fileContents struct {
	Cataloger          scope    `yaml:"cataloger" json:"cataloger" mapstructure:"cataloger"`
	SkipFilesAboveSize int64    `yaml:"skip-files-above-size" json:"skip-files-above-size" mapstructure:"skip-files-above-size"`
	Globs              []string `yaml:"globs" json:"globs" mapstructure:"globs"`
}

func defaultFileContents() fileContents {
	return fileContents{
		Cataloger: scope{
			Scope: source.SquashedScope.String(),
		},
		SkipFilesAboveSize: 1 * file.MB,
	}
}
