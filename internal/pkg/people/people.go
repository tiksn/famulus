package people

import (
	"encoding/csv"
	"errors"
	"os"
)

type People interface {
}

type people struct {
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

	return &people{
		records: records,
	}, nil
}
