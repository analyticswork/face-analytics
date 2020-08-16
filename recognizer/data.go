package recognizer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// SaveData saves data to JSON formatted file
//
func (r *Recognizer) SaveData(path string) error {
	data, err := jsonMarshal(r.Dataset)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0777)
}

// LoadData loads the data from the JSON file into the Data
//
func (r *Recognizer) LoadData(path string) error {
	if !fileExists(path) {
		return errors.New("file not found")
	}
	file, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	data := []Data{}
	if err = json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	r.Dataset = append(r.Dataset, data...)

	return nil
}
