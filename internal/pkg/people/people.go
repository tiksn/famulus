package people

import (
	"encoding/csv"
	"errors"
	"os"
)

type People interface {
	SaveToFile(path string) error
}

type people struct {
	indices map[string]int
	records [][]string
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
	for headerIndex, header := range records[0] {
		indices[header] = headerIndex
	}
	records = records[1:]
	return &people{
		indices: indices,
		records: records,
	}, nil
}

func (c *people) SaveToFile(path string) error {
	return nil
}
