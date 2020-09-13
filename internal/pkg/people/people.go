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
	indices       map[string]int
	phoneIndices  []int
	headerRecords []string
	records       [][]string
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
	headerRecords := records[0]
	records = records[1:]
	return &people{
		indices:       indices,
		phoneIndices:  phoneIndices,
		headerRecords: headerRecords,
		records:       records,
	}, nil
}

func (c *people) SaveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	allRecords := append(c.records, c.headerRecords)
	last := len(allRecords) - 1
	copy(allRecords[1:], allRecords[:last])
	allRecords[0] = c.headerRecords
	return writer.WriteAll(allRecords)
}

func (c *people) AddOrUpdate(phoneNumbers []string) error {
	i, found, err := c.findRecordIndexByNumbers(phoneNumbers)
	if err != nil {
		return err
	}

	if found {
		updatePhonenumbers(phoneNumbers, c.records[i], c.phoneIndices)
	} else {
		newRecord := make([]string, len(c.indices))

		updatePhonenumbers(phoneNumbers, newRecord, c.phoneIndices)

		c.records = append(c.records, newRecord)
	}

	return nil
}

func (c *people) findRecordIndexByNumbers(phoneNumbers []string) (int, bool, error) {
	for _, phoneNumber := range phoneNumbers {
		i, found, err := c.findRecordIndexByNumber(phoneNumber)

		if err != nil {
			return 0, false, err
		}

		if found {
			return i, true, nil
		}
	}

	return 0, false, nil
}

func (c *people) findRecordIndexByNumber(phoneNumber string) (int, bool, error) {
	for i, record := range c.records {
		for _, phoneIndex := range c.phoneIndices {
			if record[phoneIndex] == phoneNumber {
				return i, true, nil
			}
		}
	}

	return 0, false, nil
}

func updatePhonenumbers(phoneNumbers []string, record []string, phoneIndices []int) {
	for i := 0; i < len(phoneIndices) && i < len(phoneNumbers); i++ {
		record[phoneIndices[i]] = phoneNumbers[i]
	}
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
