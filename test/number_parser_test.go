package phone

import (
	"reflect"
	"testing"

	"github.com/tiksn/famulus/internal/pkg/phone"
)

func TestParse(t *testing.T) {
	testData := map[string][]string{
		"0689546321":              {"+380689546321"},
		" 0689546321 ":            {"+380689546321"},
		"068531":                  {},
		"0689576321 ; 0689576322": {"+380689576321", "+380689576322"},
	}

	for testNums, expectedNums := range testData {
		t.Logf("Parsing %s", testNums)
		nums, err := phone.Parse(testNums, "UA")
		if len(expectedNums) == 0 {
			if err == nil {
				t.Error("Should not been parsed.")
			} else {
				t.Log(err)
			}
		} else {
			if err != nil {
				t.Error(err)
			} else {
				if !reflect.DeepEqual(nums, expectedNums) {
					t.Errorf("Expected Numbers are %v", expectedNums)
				}
			}
		}

		t.Logf("Parsed Numbers are %v", nums)
	}

}
