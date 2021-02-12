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
	AddOrUpdate(phoneNumbers []string, websites []string, notes []string) error
	SaveToFile(path string) error
}

type people struct {
	indices        map[string]int
	phoneIndices   []int
	websiteIndices []int
	notesIndices   []int
	headerRecords  []string
	records        [][]string
}

type getTypeOrdinal func(t string) int

var theREs = map[string]*regexp.Regexp{
	"Phone":   regexp.MustCompile("^(?P<type_prefix>(\\w*\\'*(\\s*\\w)*))\\s*Phone\\s*(?P<number>\\d*)\\s*(?P<type_suffix>\\w*)$"),
	"Website": regexp.MustCompile("^(?P<type_prefix>(\\w*))\\s*Web Page\\s*(?P<number>\\d*)\\s*(?P<type_suffix>\\w*)$"),
	"Notes":   regexp.MustCompile("^(?P<type_prefix>)Notes(?P<number>)(?P<type_suffix>)$"),
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
		"Phone":   make(map[string]map[int]int),
		"Website": make(map[string]map[int]int),
		"Notes":   make(map[string]map[int]int),
		"Others":  make(map[string]map[int]int),
	}

	for headerIndex, header := range records[0] {
		indices[header] = headerIndex
		found := false
		for category, theRE := range theREs {
			matches := theRE.FindStringSubmatch(header)
			if matches != nil {
				t := extractType(matches, theRE)
				number := matches[theRE.SubexpIndex("number")]
				if categorizedIndices[category][t] == nil {
					categorizedIndices[category][t] = make(map[int]int)
				}

				if number == "" {
					categorizedIndices[category][t][1] = headerIndex
				} else {
					i, err := strconv.Atoi(number)
					if err != nil {
						return nil, err
					}
					categorizedIndices[category][t][i] = headerIndex
				}

				found = true
				break
			}
		}

		if !found {
			categorizedIndices["Others"][header] = map[int]int{
				1: headerIndex,
			}
		}
	}

	phoneIndices := orderPhones(categorizedIndices["Phone"])
	websiteIndices := orderWebsites(categorizedIndices["Website"])
	notesIndices := orderNotes(categorizedIndices["Notes"])
	headerRecords := records[0]
	records = records[1:]
	return &people{
		indices:        indices,
		phoneIndices:   phoneIndices,
		websiteIndices: websiteIndices,
		notesIndices:   notesIndices,
		headerRecords:  headerRecords,
		records:        records,
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

func (c *people) AddOrUpdate(phoneNumbers []string, websites []string, notes []string) error {
	i, found, err := c.findRecordIndexByValues(phoneNumbers, c.phoneIndices)
	if err != nil {
		return err
	}

	if found {
		updateValues(phoneNumbers, c.records[i], c.phoneIndices)
		updateValues(websites, c.records[i], c.websiteIndices)
		updateValues(notes, c.records[i], c.notesIndices)
	} else {
		newRecord := make([]string, len(c.indices))

		updateValues(phoneNumbers, newRecord, c.phoneIndices)
		updateValues(websites, newRecord, c.websiteIndices)
		updateValues(notes, newRecord, c.notesIndices)

		c.records = append(c.records, newRecord)
	}

	return nil
}

func (c *people) findRecordIndexByValues(values []string, indices []int) (int, bool, error) {
	for _, value := range values {
		i, found, err := c.findRecordIndexByValue(value, indices)

		if err != nil {
			return 0, false, err
		}

		if found {
			return i, true, nil
		}
	}

	return 0, false, nil
}

func (c *people) findRecordIndexByValue(value string, indices []int) (int, bool, error) {
	for i, record := range c.records {
		for _, idx := range indices {
			if record[idx] == value {
				return i, true, nil
			}
		}
	}

	return 0, false, nil
}

func updateValues(values []string, record []string, indices []int) {
	for i := 0; i < len(indices) && i < len(values); i++ {
		record[indices[i]] = values[i]
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
	return flattenIndexHierarchy(original, getPhoneTypeOrdinal)
}

func orderWebsites(original map[string]map[int]int) []int {
	return flattenIndexHierarchy(original, getWebsitesTypeOrdinal)
}

func orderNotes(original map[string]map[int]int) []int {
	return flattenIndexHierarchy(original, getNotesTypeOrdinal)
}

func flattenIndexHierarchy(original map[string]map[int]int, getOrdinal getTypeOrdinal) []int {
	keys := make([]string, 0, len(original))
	for k := range original {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return getOrdinal(keys[i]) < getOrdinal(keys[j])
	})

	var phoneIndices []int
	for _, k := range keys {
		for _, n := range original[k] {
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
	case "Company Main":
		return 5
	case "Car":
		return 6
	case "Radio":
		return 7
	case "Other":
		return 8
	case "Assistant's":
		return 9

	default:
		return 1<<31 - 1
	}
}

func getWebsitesTypeOrdinal(t string) int {
	switch t {
	case "":
		return 1
	case "Personal":
		return 2
	default:
		return 1<<31 - 1
	}

}
func getNotesTypeOrdinal(t string) int {
	switch t {
	case "":
		return 1
	default:
		return 1<<31 - 1
	}
}
