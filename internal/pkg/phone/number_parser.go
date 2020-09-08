package phone

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func Parse(numbers string, defaultRegion string) ([]string, error) {
	nums, err := splitNumbers(numbers)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, num := range nums {
		pnum, err := phonenumbers.Parse(num, defaultRegion)
		if err != nil {
			return nil, err
		}

		snum := phonenumbers.Format(pnum, phonenumbers.E164)

		if !phonenumbers.IsValidNumber(pnum) {
			return nil, fmt.Errorf("%s Is Not Valid Number", snum)
		}

		result = append(result, snum)
	}

	return result, nil
}

func splitNumbers(numbers string) ([]string, error) {
	nums := strings.Split(numbers, ";")
	return nums, nil
	// return nil, fmt.Errorf("Unable to split numbers %s", numbers)
}
