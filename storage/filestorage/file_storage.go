package filestorage

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"program/model"
)

type FileStorage struct {
	fileName string
}

func NewFileStorage(file string) *FileStorage {

	return &FileStorage{
		fileName: file,
	}

}

func (fs *FileStorage) Load() ([]model.Joke, error) {

	file, err := os.Open(fs.fileName)

	if err != nil {
		return nil, errors.Wrap(err, "unable to open a file")
	}
	defer file.Close()

	var result []model.Joke
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding file")
	}
	return result, nil

}

func (fs *FileStorage) Save(jokes []model.Joke) error {
	structBytes, err := json.MarshalIndent(jokes, "", " ")
	if err != nil {
		return fmt.Errorf(" error marshalling JSON:%w:", err)
	}
	err = ioutil.WriteFile("db/reddit_jokes.json", structBytes, 0644)
	if err != nil {
		return fmt.Errorf(" error saving file:%w", err)
	}

	return nil
}
