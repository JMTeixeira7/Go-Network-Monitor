package storage

import (
	"encoding/json"
	"os"
)

type Stats_Storage interface {
	Save([]Search) error
	LoadAll() ([]Search, error)
	Load() (Search, error)
}

type FileStorage struct {
	Filename string
}

func (f FileStorage) LoadAll() ([]Search, error) {
	data, err := os.ReadFile(f.Filename)

	if (err != nil) {
		if os.IsNotExist(err) { //file doesnt exist
			return []Search{}, nil
		}
		return nil, err
	}

	var searches []Search
	err = json.Unmarshal(data, &searches)
	if err != nil {
		return nil, err
	}
	return searches, nil
}

func (f FileStorage) Save(input []Search) error {
	data, err := os.ReadFile(f.Filename)

	if (err != nil) {
		return err
	}

	var searches []Search
	err = json.Unmarshal(data, &searches)

	searches = append(searches, input...)
	data, err = json.MarshalIndent(searches, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.Filename, data, 0644)
}