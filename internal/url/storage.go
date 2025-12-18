package url

import (
	"encoding/json"
	"os"
)

type URLStorage interface{
	Load() ([]Url, error)
	Save([]Url) error
	Set() error
}

type FileStorage struct {
	Filename string
}

func (f FileStorage) Load() ([]Url, error) {
	data, err := os.ReadFile(f.Filename)

	if (err != nil) {
		if os.IsNotExist(err) {	// file doesnt exist
			return []Url{}, nil
		}
		return nil, err
	}

	var urls []Url
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (f FileStorage) Save(urls []Url) error {
	data, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.Filename, data, 0644) //premission R/W, user R only
}

func (f FileStorage) Set() error {
	SetTargetURLs()
	data, err := json.MarshalIndent(GetTargetURLs(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.Filename, data, 0644) //premission R/W, user R only
}