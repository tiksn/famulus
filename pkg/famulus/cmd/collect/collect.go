package collect

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	config "github.com/tiksn/famulus/internal/app/famulus"
	"github.com/tiksn/famulus/internal/pkg/people"
	"github.com/tiksn/famulus/internal/pkg/scraper"
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
		adr, err := c.GetAddress(args[0])
		if err != nil {
			return err
		}

		pageNumber := 1

		if argCount == 2 {
			pageNumberParsed, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			pageNumber = pageNumberParsed
		}

		csv, err := c.GetContactsCsvFilePath()
		if err != nil {
			return err
		}

		peopleDB, err := people.LoadFromFile(csv)
		if err != nil {
			return err
		}

		adrUrl, err := adr.GetAddress()
		if err != nil {
			return err
		}
		adrUrl = strings.ReplaceAll(adrUrl, "{page_number}", strconv.Itoa(pageNumber))
		fmt.Println(adr)

		phonrUrl, err := adr.GetPhoneAddress()
		if err != nil {
			return err
		}

		contacts, err := scraper.ListScrape(adrUrl, phonrUrl)
		if err != nil {
			return err
		}

		for _, contact := range contacts {
			fmt.Println(contact)
		}

		return peopleDB.SaveToFile(csv)
	}

	return nil
}
