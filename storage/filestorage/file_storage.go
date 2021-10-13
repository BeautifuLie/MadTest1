package filestorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"program/joker"
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

	err = json.NewDecoder(file).Decode(&joker.S.JokesStruct)
	if err != nil {
		return nil, err
	}

	for _, j := range joker.S.JokesStruct {
		joker.S.JokesMap[j.ID] = j
	}

	return joker.S.JokesStruct, nil

}

func (fs *FileStorage) Save([]model.Joke) error {
	structBytes, err := json.MarshalIndent(joker.S.JokesStruct, "", " ")
	if err != nil {
		errors.New(" Error marshalling JSON")
	}
	err = ioutil.WriteFile("reddit_jokes.json", structBytes, 0644)
	if err != nil {
		return errors.New(" Error saving file")
	}

	return nil
}
