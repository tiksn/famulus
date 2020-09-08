package phone

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

func Parse(numbers string, defaultRegion string) ([]string, error) {
	nums, err := splitNumbers(numbers)
	if err != nil {
		return nil, err
	}

	for _, num := range nums {
		pnum, err := phonenumbers.Parse(num, defaultRegion)
		if err != nil {
			return nil, err
		}

		_ = pnum
	}

	return nil, fmt.Errorf("Unable to parse numbers %s", numbers)
}

func splitNumbers(numbers string) ([]string, error) {
	return nil, fmt.Errorf("Unable to split numbers %s", numbers)
}
