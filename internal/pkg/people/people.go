package people

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type People interface {
	AddOrUpdate(phoneNumbers []string) error
	SaveToFile(path string) error
}

type people struct {
	indices      map[string]int
	phoneIndices []int
	records      [][]string
}

func LoadFromFile(path string) (People, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("No headers found")
	}

	indices := make(map[string]int)
	categorizedIndices := map[string]map[string]map[int]int{
		"Phone":  make(map[string]map[int]int),
		"Others": make(map[string]map[int]int),
	}
	phoneRE := regexp.MustCompile("^(?P<type_prefix>(\\w*\\'*\\w*))\\s*Phone\\s*(?P<number>\\d*)\\s*(?P<type_suffix>\\w*)$")
	for headerIndex, header := range records[0] {
		indices[header] = headerIndex

		matches := phoneRE.FindStringSubmatch(header)
		if matches != nil {
			t := extractType(matches, phoneRE)
			number := matches[phoneRE.SubexpIndex("number")]
			if categorizedIndices["Phone"][t] == nil {
				categorizedIndices["Phone"][t] = make(map[int]int)
			}

			if number == "" {
				categorizedIndices["Phone"][t][1] = headerIndex
			} else {
				i, err := strconv.Atoi(number)
				if err != nil {
					return nil, err
				}
				categorizedIndices["Phone"][t][i] = headerIndex
			}
		} else {
			categorizedIndices["Others"][header] = map[int]int{
				1: headerIndex,
			}
		}
	}
	phoneIndices := orderPhones(categorizedIndices["Phone"])
	records = records[1:]
	return &people{
		indices:      indices,
		phoneIndices: phoneIndices,
		records:      records,
	}, nil
}

func (c *people) SaveToFile(path string) error {
	return nil
}

func (c *people) AddOrUpdate(phoneNumbers []string) error {
	return nil
}

func extractType(matches []string, re *regexp.Regexp) string {
	typePrefix := matches[re.SubexpIndex("type_prefix")]
	typeSuffix := matches[re.SubexpIndex("type_suffix")]

	if typePrefix == "" {
		return typeSuffix
	}

	if typeSuffix == "" {
		return typePrefix
	}

	return fmt.Sprintf("%s %s", typePrefix, typeSuffix)
}

func orderPhones(original map[string]map[int]int) []int {
	keys := make([]string, 0, len(original))
	for k := range original {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return getPhoneTypeOrdinal(keys[i]) < getPhoneTypeOrdinal(keys[j])
	})

	var phoneIndices []int
	for _, v := range original {
		for _, n := range v {
			phoneIndices = append(phoneIndices, n)
		}
	}
	return phoneIndices
}

func getPhoneTypeOrdinal(t string) int {

	switch t {
	case "Mobile":
		return 1
	case "Primary":
		return 2
	case "Home":
		return 3
	case "Business":
		return 4
	case "Car":
		return 5
	case "Radio":
		return 6
	case "Other":
		return 7
	case "Assistant's":
		return 8

	default:
		return 1<<31 - 1
	}
}
