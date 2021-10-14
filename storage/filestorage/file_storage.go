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

func NewFileStorage(file string) *FileStorage {

	return &FileStorage{
		fileName: file}

}

//var f = &FileStorage{
//	fileName: "reddit_jokes.json",
//}

func (fs *FileStorage) Load() ([]model.Joke, error) {

	fs.fileName = "db/reddit_jokes1.json"

	file, err := os.Open(fs.fileName)

	if err != nil {

		return nil, fmt.Errorf(" Can't open file : %w", err)
	}

	defer file.Close()

	var result []model.Joke

	err1 := json.NewDecoder(file).Decode(&result)
	if err1 != nil {
		return nil, fmt.Errorf(" Can't decode file : %w", err1)
	}

	return result, nil

}

func (fs *FileStorage) Save(jokes []model.Joke) error {
	structBytes, err := json.MarshalIndent(jokes, "", " ")
	if err != nil {
		return fmt.Errorf(" error marshalling JSON:%w:", err)
	}
	err = ioutil.WriteFile("reddit_jokes.json", structBytes, 0644)
	if err != nil {
		return fmt.Errorf(" error saving file:%w", err)
	}

	return nil
}
