package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/anchore/clio"
	"github.com/anchore/stereoscope"
	"github.com/anubhav06/sanrakshya-cli/cmd/sanrakshya/cli/commands"
	handler "github.com/anubhav06/sanrakshya-cli/cmd/sanrakshya/cli/ui"
	"github.com/anubhav06/sanrakshya-cli/cmd/sanrakshya/internal/ui"
	"github.com/anubhav06/sanrakshya-cli/internal/bus"
	"github.com/anubhav06/sanrakshya-cli/internal/log"
	"github.com/anubhav06/sanrakshya-cli/internal/redact"
)

// Application constructs the `sanrakshya packages` command and aliases the root command to `sanrakshya packages`.
// It is also responsible for organizing flag usage and injecting the application config for each command.
// It also constructs the sanrakshya attest command and the sanrakshya version command.
// `RunE` is the earliest that the complete application configuration can be loaded.
func Application(id clio.Identification) clio.Application {
	app, _ := create(id, os.Stdout)
	return app
}

// Command returns the root command for the sanrakshya CLI application. This is useful for embedding the entire sanrakshya CLI
// into an existing application.
func Command(id clio.Identification) *cobra.Command {
	_, cmd := create(id, os.Stdout)
	return cmd
}

func create(id clio.Identification, out io.Writer) (clio.Application, *cobra.Command) {
	clioCfg := clio.NewSetupConfig(id).
		WithUIConstructor(
			// select a UI based on the logging configuration and state of stdin (if stdin is a tty)
			func(cfg clio.Config) ([]clio.UI, error) {
				noUI := ui.None(out, cfg.Log.Quiet)
				if !cfg.Log.AllowUI(os.Stdin) || cfg.Log.Quiet {
					return []clio.UI{noUI}, nil
				}

				return []clio.UI{
					ui.New(out, cfg.Log.Quiet,
						handler.New(handler.DefaultHandlerConfig()),
					),
					noUI,
				}, nil
			},
		).
		WithInitializers(
			func(state *clio.State) error {
				// clio is setting up and providing the bus, redact store, and logger to the application. Once loaded,
				// we can hoist them into the internal packages for global use.
				stereoscope.SetBus(state.Bus)
				bus.Set(state.Bus)

				redact.Set(state.RedactStore)

				log.Set(state.Logger)
				stereoscope.SetLogger(state.Logger)

				return nil
			},
		).
		WithPostRuns(func(state *clio.State, err error) {
			stereoscope.Cleanup()
		})

	app := clio.New(*clioCfg)

	// since root is aliased as the packages cmd we need to construct this command first
	// we also need the command to have information about the `root` options because of this alias
	packagesCmd := commands.Scan(app)

	// rootCmd is currently an alias for the packages command
	rootCmd := commands.Root(app, packagesCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// add sub-commands
	rootCmd.AddCommand(
		packagesCmd,
		commands.Submit(app),
		commands.Login(app),
	)

	// explicitly set Cobra output to the real stdout to write things like errors and help
	rootCmd.SetOut(out)

	return app, rootCmd
}
