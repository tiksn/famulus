package collect

import (
	"fmt"

	"github.com/spf13/cobra"
	config "github.com/tiksn/famulus/internal/app/famulus"
)

func NewCollectCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect Contacts",
		Long:  "Collect Contacts",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Collecting ....")
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}
