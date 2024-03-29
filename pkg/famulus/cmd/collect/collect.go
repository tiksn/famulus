package collect

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	config "github.com/tiksn/famulus/internal/app/famulus"
	"github.com/tiksn/famulus/internal/pkg/people"
	"github.com/tiksn/famulus/internal/pkg/phone"
	"github.com/tiksn/famulus/internal/pkg/scraper"
)

func NewCollectCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect <type> [page_number]",
		Short: "Collect Contacts",
		Long:  "Collect Contacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, err := cmd.Flags().GetDuration("interval")
			if err != nil {
				return err
			}
			return collectCmd(*c, args, interval)
		},
		Args: cobra.MaximumNArgs(2),
	}

	cmd.Flags().DurationP("interval", "i", 5*time.Second, "Delay before making HTTP request")
	return cmd
}

func collectCmd(c config.Config, args []string, interval time.Duration) error {
	argCount := len(args)

	if argCount == 0 {
		a, err := c.ListSources()
		if err != nil {
			return err
		}
		for _, n := range a {
			fmt.Println("Available collection option:", n)
		}
		return nil
	} else if argCount == 1 || argCount == 2 {
		src, err := c.GetSource(args[0])
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

		kind, err := src.GetKind()
		if err != nil {
			return err
		}

		srcUrl, err := src.GetAddress()
		if err != nil {
			return err
		}
		srcUrl = strings.ReplaceAll(srcUrl, "{page_number}", strconv.Itoa(pageNumber))

		phoneUrl, err := src.GetPhoneAddress()
		if err != nil {
			return err
		}

		kindScraper, err := scraper.GetScraper(kind, srcUrl, phoneUrl, interval)
		if err != nil {
			return err
		}

		contacts, err := kindScraper.ListScrape()
		if err != nil {
			return err
		}

		for _, contact := range contacts {
			region, err := src.GetDefaultRegion()
			if err != nil {
				return err
			}
			numbers, err := phone.Parse(contact.GetNumbers(), region)
			if err != nil {
				return err
			}

			notes := contact.GetTitle() + " --- " + contact.GetDescription()
			err = peopleDB.AddOrUpdate(numbers, []string{contact.GetWebsite()}, []string{notes})
			if err != nil {
				return err
			}
		}

		err = peopleDB.SaveToFile(csv)
		if err != nil {
			return err
		}
	}

	return nil
}
