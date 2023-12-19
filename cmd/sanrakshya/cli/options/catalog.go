package options

import (
	"fmt"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/go-homedir"
	"github.com/scylladb/go-set/strset"

	"github.com/anchore/clio"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/cataloging"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger"
	golangCataloger "github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/golang"
	javaCataloger "github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/java"
	javascriptCataloger "github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/javascript"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/kernel"
	pythonCataloger "github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/python"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/source"
)

type Catalog struct {
	Package                         pkg          `yaml:"package" json:"package" mapstructure:"package"`
	Golang                          golang       `yaml:"golang" json:"golang" mapstructure:"golang"`
	Java                            java         `yaml:"java" json:"java" mapstructure:"java"`
	Javascript                      javascript   `yaml:"javascript" json:"javascript" mapstructure:"javascript"`
	LinuxKernel                     linuxKernel  `yaml:"linux-kernel" json:"linux-kernel" mapstructure:"linux-kernel"`
	Python                          python       `yaml:"python" json:"python" mapstructure:"python"`
	FileMetadata                    fileMetadata `yaml:"file-metadata" json:"file-metadata" mapstructure:"file-metadata"`
	FileContents                    fileContents `yaml:"file-contents" json:"file-contents" mapstructure:"file-contents"`
	Registry                        registry     `yaml:"registry" json:"registry" mapstructure:"registry"`
	Name                            string       `yaml:"name" json:"name" mapstructure:"name"`
	Source                          sourceCfg    `yaml:"source" json:"source" mapstructure:"source"`
	Parallelism                     int          `yaml:"parallelism" json:"parallelism" mapstructure:"parallelism"`                                                                         // the number of catalog workers to run in parallel
	DefaultImagePullSource          string       `yaml:"default-image-pull-source" json:"default-image-pull-source" mapstructure:"default-image-pull-source"`                               // specify default image pull source
	ExcludeBinaryOverlapByOwnership bool         `yaml:"exclude-binary-overlap-by-ownership" json:"exclude-binary-overlap-by-ownership" mapstructure:"exclude-binary-overlap-by-ownership"` // exclude synthetic binary packages owned by os package files
}

var _ interface {
	clio.FlagAdder
	clio.PostLoader
} = (*Catalog)(nil)

func DefaultCatalog() Catalog {
	return Catalog{
		Package:                         defaultPkg(),
		LinuxKernel:                     defaultLinuxKernel(),
		FileMetadata:                    defaultFileMetadata(),
		FileContents:                    defaultFileContents(),
		Source:                          defaultSourceCfg(),
		Parallelism:                     1,
		ExcludeBinaryOverlapByOwnership: true,
	}
}

func (cfg *Catalog) AddFlags(flags clio.FlagSet) {
	var validScopeValues []string
	for _, scope := range source.AllScopes {
		validScopeValues = append(validScopeValues, strcase.ToDelimited(string(scope), '-'))
	}
	flags.StringVarP(&cfg.Package.Cataloger.Scope, "scope", "s",
		fmt.Sprintf("selection of layers to catalog, options=%v", validScopeValues))

}

func (cfg *Catalog) PostLoad() error {
	// parse options on this struct
	var catalogers []string
	sort.Strings(catalogers)

	if err := checkDefaultSourceValues(cfg.DefaultImagePullSource); err != nil {
		return err
	}

	return nil
}

func (cfg Catalog) ToCatalogerConfig() cataloger.Config {
	return cataloger.Config{
		Search: cataloger.SearchConfig{
			IncludeIndexedArchives:   cfg.Package.SearchIndexedArchives,
			IncludeUnindexedArchives: cfg.Package.SearchUnindexedArchives,
			Scope:                    cfg.Package.Cataloger.GetScope(),
		},
		Parallelism: cfg.Parallelism,
		Golang: golangCataloger.DefaultCatalogerConfig().
			WithSearchLocalModCacheLicenses(cfg.Golang.SearchLocalModCacheLicenses).
			WithLocalModCacheDir(cfg.Golang.LocalModCacheDir).
			WithSearchRemoteLicenses(cfg.Golang.SearchRemoteLicenses).
			WithProxy(cfg.Golang.Proxy).
			WithNoProxy(cfg.Golang.NoProxy),
		LinuxKernel: kernel.LinuxKernelCatalogerConfig{
			CatalogModules: cfg.LinuxKernel.CatalogModules,
		},
		Java: javaCataloger.DefaultArchiveCatalogerConfig().
			WithUseNetwork(cfg.Java.UseNetwork).
			WithMavenBaseURL(cfg.Java.MavenURL).
			WithArchiveTraversal(
				cataloging.ArchiveSearchConfig{
					IncludeIndexedArchives:   cfg.Package.SearchIndexedArchives,
					IncludeUnindexedArchives: cfg.Package.SearchUnindexedArchives,
				},
				cfg.Java.MaxParentRecursiveDepth),
		Javascript: javascriptCataloger.DefaultCatalogerConfig().
			WithSearchRemoteLicenses(cfg.Javascript.SearchRemoteLicenses).
			WithNpmBaseURL(cfg.Javascript.NpmBaseURL),
		Python: pythonCataloger.CatalogerConfig{
			GuessUnpinnedRequirements: cfg.Python.GuessUnpinnedRequirements,
		},
		ExcludeBinaryOverlapByOwnership: cfg.ExcludeBinaryOverlapByOwnership,
	}
}

var validDefaultSourceValues = []string{"registry", "docker", "podman", ""}

func checkDefaultSourceValues(source string) error {
	validValues := strset.New(validDefaultSourceValues...)
	if !validValues.Has(source) {
		validValuesString := strings.Join(validDefaultSourceValues, ", ")
		return fmt.Errorf("%s is not a valid default source; please use one of the following: %s''", source, validValuesString)
	}

	return nil
}

func expandFilePath(file string) (string, error) {
	if file != "" {
		expandedPath, err := homedir.Expand(file)
		if err != nil {
			return "", fmt.Errorf("unable to expand file path=%q: %w", file, err)
		}
		file = expandedPath
	}
	return file, nil
}
