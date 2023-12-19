package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anchore/clio"
)

const (
	rootLong = `
	Sanrakshya is a CLI tool for generating SBOMs from container images and filesystems.
	Made with 💖 by Team PiedPiper & developed indigeniously in India.
	`
)

func Root(app clio.Application, packagesCmd *cobra.Command) *cobra.Command {
	id := app.ID()

	opts := defaultScanOptions()

	return app.SetupRootCommand(&cobra.Command{
		Use:   fmt.Sprintf("%s [SOURCE]", app.ID().Name),
		Short: "Sanrakshya is a CLI tool for generating SBOMs from container images and filesystems",
		Long:  rootLong,
		Args:  packagesCmd.Args,
		RunE: func(cmd *cobra.Command, args []string) error {
			// restoreStdout := ui.CaptureStdoutToTraceLog()
			// defer restoreStdout()

			return runScan(id, opts, args[0])
		},
	}, opts)
}
