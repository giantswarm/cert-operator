package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	CLICmd = &cobra.Command{
		Use:   "certctl",
		Short: "A command line tool able to request certificate generation from Vault to write certificate files to the local filesystem.",

		Run: cliRun,
	}
)

func cliRun(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
	os.Exit(1)
}
