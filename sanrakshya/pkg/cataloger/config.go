package cataloger

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/cataloging"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/golang"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/java"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/javascript"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/kernel"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/python"
)

// TODO: these field naming vs helper function naming schemes are inconsistent.
type Config struct {
	Search                          SearchConfig
	Golang                          golang.CatalogerConfig
	LinuxKernel                     kernel.LinuxKernelCatalogerConfig
	Python                          python.CatalogerConfig
	Java                            java.ArchiveCatalogerConfig
	Javascript                      javascript.CatalogerConfig
	Catalogers                      []string
	Parallelism                     int
	ExcludeBinaryOverlapByOwnership bool
}

func DefaultConfig() Config {
	return Config{
		Search:                          DefaultSearchConfig(),
		Parallelism:                     1,
		LinuxKernel:                     kernel.DefaultLinuxCatalogerConfig(),
		Python:                          python.DefaultCatalogerConfig(),
		Java:                            java.DefaultArchiveCatalogerConfig(),
		Javascript:                      javascript.DefaultCatalogerConfig(),
		ExcludeBinaryOverlapByOwnership: true,
	}
}

// JavaConfig merges relevant config values from Config to return a java.Config struct.
// Values like IncludeUnindexedArchives and IncludeIndexedArchives are used across catalogers
// and are not specific to Java requiring this merge.
func (c Config) JavaConfig() java.ArchiveCatalogerConfig {
	return java.ArchiveCatalogerConfig{
		ArchiveSearchConfig: cataloging.ArchiveSearchConfig{
			IncludeUnindexedArchives: c.Search.IncludeUnindexedArchives,
			IncludeIndexedArchives:   c.Search.IncludeIndexedArchives,
		},
		UseNetwork:              c.Java.UseNetwork,
		MavenBaseURL:            c.Java.MavenBaseURL,
		MaxParentRecursiveDepth: c.Java.MaxParentRecursiveDepth,
	}
}
