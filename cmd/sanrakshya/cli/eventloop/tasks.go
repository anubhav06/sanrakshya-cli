package eventloop

import (
	"github.com/anubhav06/sanrakshya-cli/cmd/sanrakshya/cli/options"
	"github.com/anubhav06/sanrakshya-cli/internal/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/artifact"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file/cataloger/filecontent"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file/cataloger/filedigest"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file/cataloger/filemetadata"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/source"
)

type Task func(*sbom.Artifacts, source.Source) ([]artifact.Relationship, error)

func Tasks(opts *options.Catalog) ([]Task, error) {
	var tasks []Task

	generators := []func(opts *options.Catalog) (Task, error){
		generateCatalogPackagesTask,
		generateCatalogFileMetadataTask,
		generateCatalogFileDigestsTask,
		generateCatalogContentsTask,
	}

	for _, generator := range generators {
		task, err := generator(opts)
		if err != nil {
			return nil, err
		}

		if task != nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func generateCatalogPackagesTask(opts *options.Catalog) (Task, error) {
	if !opts.Package.Cataloger.Enabled {
		return nil, nil
	}

	task := func(results *sbom.Artifacts, src source.Source) ([]artifact.Relationship, error) {
		packageCatalog, relationships, theDistro, err := sanrakshya.CatalogPackages(src, opts.ToCatalogerConfig())

		results.Packages = packageCatalog
		results.LinuxDistribution = theDistro

		return relationships, err
	}

	return task, nil
}

func generateCatalogFileMetadataTask(opts *options.Catalog) (Task, error) {
	if !opts.FileMetadata.Cataloger.Enabled {
		return nil, nil
	}

	metadataCataloger := filemetadata.NewCataloger()

	task := func(results *sbom.Artifacts, src source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(opts.FileMetadata.Cataloger.GetScope())
		if err != nil {
			return nil, err
		}

		result, err := metadataCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileMetadata = result
		return nil, nil
	}

	return task, nil
}

func generateCatalogFileDigestsTask(opts *options.Catalog) (Task, error) {
	if !opts.FileMetadata.Cataloger.Enabled {
		return nil, nil
	}

	hashes, err := file.Hashers(opts.FileMetadata.Digests...)
	if err != nil {
		return nil, err
	}

	digestsCataloger := filedigest.NewCataloger(hashes)

	task := func(results *sbom.Artifacts, src source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(opts.FileMetadata.Cataloger.GetScope())
		if err != nil {
			return nil, err
		}

		result, err := digestsCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileDigests = result
		return nil, nil
	}

	return task, nil
}

func generateCatalogContentsTask(opts *options.Catalog) (Task, error) {
	if !opts.FileContents.Cataloger.Enabled {
		return nil, nil
	}

	contentsCataloger, err := filecontent.NewCataloger(opts.FileContents.Globs, opts.FileContents.SkipFilesAboveSize) //nolint:staticcheck
	if err != nil {
		return nil, err
	}

	task := func(results *sbom.Artifacts, src source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(opts.FileContents.Cataloger.GetScope())
		if err != nil {
			return nil, err
		}

		result, err := contentsCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileContents = result
		return nil, nil
	}

	return task, nil
}

func RunTask(t Task, a *sbom.Artifacts, src source.Source, c chan<- artifact.Relationship) error {
	defer close(c)

	relationships, err := t(a, src)
	if err != nil {
		return err
	}

	for _, relationship := range relationships {
		c <- relationship
	}

	return nil
}
