package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type FileStorage struct {
	baseDir  string
	fileName string
}

var F = FileStorage{
	baseDir:  "C:\\Users\\Lafazan\\GO\\MadAppGangTest",
	fileName: "reddit_jokes.json",
}

func (fs *FileStorage) Load() ([]Joke, error) {

	file, err := os.Open(F.fileName)
	if err != nil {
		return nil, fmt.Errorf("Can't open file %s: %w", F.fileName, err)
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&S.JokesStruct)
	if err != nil {
		return nil, err
	}

	for _, j := range S.JokesStruct {
		S.JokesMap[j.ID] = j
	}

	return S.JokesStruct, nil

}

func (fs *FileStorage) Save([]Joke) error {
	structBytes, err := json.MarshalIndent(S.JokesStruct, "", " ")
	if err != nil {
		errors.New(" Error marshalling JSON")
	}
	err = ioutil.WriteFile("reddit_jokes.json", structBytes, 0644)
	if err != nil {
		return errors.New(" Error saving file")
	}

	return nil
}
