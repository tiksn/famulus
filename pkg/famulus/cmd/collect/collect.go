package collect

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	config "github.com/tiksn/famulus/internal/app/famulus"
)

func NewCollectCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect <type> [page_number]",
		Short: "Collect Contacts",
		Long:  "Collect Contacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectCmd(*c, args)
		},
		Args: cobra.MaximumNArgs(2),
	}

	return cmd
}

func collectCmd(c config.Config, args []string) error {
	argCount := len(args)

	if argCount == 0 {
		a, err := c.ListAddress()
		if err != nil {
			return err
		}
		for _, n := range a {
			fmt.Println("Available collection option:", n)
		}
		return nil
	} else if argCount == 1 || argCount == 2 {
		adr, err2 := c.GetAddress(args[0])
		if err2 != nil {
			return err2
		}

		pageNumber := 1

		if argCount == 2 {
			pageNumberParsed, err3 := strconv.Atoi(args[1])
			if err3 != nil {
				return err3
			}
			pageNumber = pageNumberParsed
		}

		adr = strings.ReplaceAll(adr, "{page_number}", strconv.Itoa(pageNumber))
		fmt.Println(adr)
		csv, err3 := c.GetContactsCsvFilePath()
		if err3 != nil {
			return err3
		}
		fmt.Println(csv)
	}

	return nil
}
