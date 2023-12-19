package commands

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/anchore/clio"
	"github.com/anubhav06/sanrakshya-cli/cmd/sanrakshya/cli/eventloop"
	"github.com/anubhav06/sanrakshya-cli/cmd/sanrakshya/cli/options"
	"github.com/anubhav06/sanrakshya-cli/internal"
	"github.com/anubhav06/sanrakshya-cli/internal/bus"
	"github.com/anubhav06/sanrakshya-cli/internal/file"
	"github.com/anubhav06/sanrakshya-cli/internal/log"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/artifact"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/sbom"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/source"
)

const (
	scanExample = `  {{.appName}} {{.command}} alpine:latest                                a summary of discovered packages
  {{.appName}} {{.command}} alpine:latest -o spdx-json                   show a SPDX 2.3 JSON formatted SBOM

  Supports the following image sources:
    {{.appName}} {{.command}} yourrepo/yourimage:tag     defaults to using images from a Docker daemon, otherwise registry.
    {{.appName}} {{.command}} path/to/a/file/or/dir      any local filesystem path (directory or file)
`
	scanHelp = scanExample
)

type scanOptions struct {
	options.Config      `yaml:",inline" mapstructure:",squash"`
	options.Output      `yaml:",inline" mapstructure:",squash"`
	options.UpdateCheck `yaml:",inline" mapstructure:",squash"`
	options.Catalog     `yaml:",inline" mapstructure:",squash"`
}

func defaultScanOptions() *scanOptions {
	return &scanOptions{
		Output:      options.DefaultOutput(),
		UpdateCheck: options.DefaultUpdateCheck(),
		Catalog:     options.DefaultCatalog(),
	}
}

//nolint:dupl
func Scan(app clio.Application) *cobra.Command {
	id := app.ID()

	opts := defaultScanOptions()

	return app.SetupCommand(&cobra.Command{
		Use:   "scan [SOURCE]",
		Short: "Scan the source to generate a SBOM",
		Long:  "Generate a packaged-based Software Bill Of Materials (SBOM) from container images and filesystems",
		Example: internal.Tprintf(scanHelp, map[string]interface{}{
			"appName": id.Name,
			"command": "scan",
		}),
		Args: validateScanArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// restoreStdout := ui.CaptureStdoutToTraceLog()
			// defer restoreStdout()

			return runScan(id, opts, args[0])
		},
	}, opts)
}

func validateScanArgs(cmd *cobra.Command, args []string) error {
	return validateArgs(cmd, args, "an image/directory argument is required")
}

func validateArgs(cmd *cobra.Command, args []string, error string) error {
	if len(args) == 0 {
		// in the case that no arguments are given we want to show the help text and return with a non-0 return code.
		if err := cmd.Help(); err != nil {
			return fmt.Errorf("unable to display help: %w", err)
		}
		return fmt.Errorf(error)
	}

	return cobra.MaximumNArgs(1)(cmd, args)
}

// nolint:funlen
func runScan(id clio.Identification, opts *scanOptions, userInput string) error {
	writer, err := opts.SBOMWriter()
	if err != nil {
		return err
	}

	src, err := getSource(&opts.Catalog, userInput)

	if err != nil {
		return err
	}

	defer func() {
		if src != nil {
			if err := src.Close(); err != nil {
				log.Tracef("unable to close source: %+v", err)
			}
		}
	}()

	s, err := generateSBOM(id, src, &opts.Catalog)
	if err != nil {
		return err
	}

	if s == nil {
		return fmt.Errorf("no SBOM produced for %q", userInput)
	}

	if err := writer.Write(*s); err != nil {
		return fmt.Errorf("failed to write SBOM: %w", err)
	}

	return nil
}

func getSource(opts *options.Catalog, userInput string, filters ...func(*source.Detection) error) (source.Source, error) {
	detection, err := source.Detect(
		userInput,
		source.DetectConfig{
			DefaultImageSource: opts.DefaultImagePullSource,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not deteremine source: %w", err)
	}

	for _, filter := range filters {
		if err := filter(detection); err != nil {
			return nil, err
		}
	}

	hashers, err := file.Hashers(opts.Source.File.Digests...)
	if err != nil {
		return nil, fmt.Errorf("invalid hash: %w", err)
	}

	src, err := detection.NewSource(
		source.DetectionSourceConfig{
			RegistryOptions:  opts.Registry.ToOptions(),
			DigestAlgorithms: hashers,
		},
	)

	if err != nil {
		if userInput == "power-user" {
			bus.Notify("Note: the 'power-user' command has been removed.")
		}
		return nil, fmt.Errorf("failed to construct source from user input %q: %w", userInput, err)
	}

	return src, nil
}

func generateSBOM(id clio.Identification, src source.Source, opts *options.Catalog) (*sbom.SBOM, error) {
	tasks, err := eventloop.Tasks(opts)
	if err != nil {
		return nil, err
	}

	s := sbom.SBOM{
		Source: src.Describe(),
		Descriptor: sbom.Descriptor{
			Name:          id.Name,
			Version:       id.Version,
			Configuration: opts,
		},
	}

	err = buildRelationships(&s, src, tasks)

	return &s, err
}

func buildRelationships(s *sbom.SBOM, src source.Source, tasks []eventloop.Task) error {
	var errs error

	var relationships []<-chan artifact.Relationship
	for _, task := range tasks {
		c := make(chan artifact.Relationship)
		relationships = append(relationships, c)
		go func(task eventloop.Task) {
			err := eventloop.RunTask(task, &s.Artifacts, src, c)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}(task)
	}

	s.Relationships = append(s.Relationships, mergeRelationships(relationships...)...)

	return errs
}

func mergeRelationships(cs ...<-chan artifact.Relationship) (relationships []artifact.Relationship) {
	for _, c := range cs {
		for n := range c {
			relationships = append(relationships, n)
		}
	}

	return relationships
}
