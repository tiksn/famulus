package root

import (
	"fmt"

	"github.com/spf13/cobra"
	config "github.com/tiksn/famulus/internal/app/famulus"
	"github.com/tiksn/famulus/pkg/famulus/cmd/collect"
)

const (
	AppVersion = "1.0.0"
)

func NewCmdRoot(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "famulus <command> <subcommand> [flags]",
		Short: "famulus CLI",
		Long:  `Collect contacts from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,

		Version: AppVersion,
	}

	versionOutput := fmt.Sprintf("famulus version %s\n", cmd.Version)
	cmd.AddCommand(&cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(versionOutput)
		},
	})

	cmd.AddCommand(collect.NewCollectCmd(c))

	return cmd
}
