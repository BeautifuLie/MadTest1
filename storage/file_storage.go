package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//type FileStorage struct {
//	baseDir string
//}
//
//var f = FileStorage{
//	baseDir: "C:/Users/Lafazan/go/src/MadAppGangTest",
//}

func (fs *Server) Load() ([]Joke, error) {
	fn := "reddit_jokes.json"

	file, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("Can't open file %s: %w", fn, err)
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

func (fs *Server) Save([]Joke) error {
	structBytes, err := json.MarshalIndent(S.JokesStruct, "", " ")
	if err != nil {
		fmt.Println("Error marshalling JSON", err)
	}
	err = ioutil.WriteFile("jokesStruct.json", structBytes, 0644)

	return err
}
