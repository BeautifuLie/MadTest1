package filestorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"program/model"
)

type FileStorage struct {
	fileName string
}

func NewFileStorage(fileName string) *FileStorage {
	return f
	//return &FileStorage{
	//
	//	fileName: fileName}

}

var f = &FileStorage{
	fileName: "reddit_jokes.json",
}

func (fs *FileStorage) Load() ([]model.Joke, error) {

	file, err := os.Open(f.fileName)

	if err != nil {
		return nil, fmt.Errorf("Can't open file %s: %w", fs.fileName, err)
	}

	defer file.Close()

	var result []model.Joke

	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (fs *FileStorage) Save(jokes []model.Joke) error {
	structBytes, err := json.MarshalIndent(jokes, "", " ")
	if err != nil {
		fmt.Errorf(" error marshalling JSON:%w:", err)
	}
	err = ioutil.WriteFile("reddit_jokes2.json", structBytes, 0644)
	if err != nil {
		return fmt.Errorf(" error saving file:%w", err)
	}

	return nil
}
