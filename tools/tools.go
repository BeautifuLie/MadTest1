package tools

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"program/model"
	"time"
)

func RandTimeAndUserID() (time.Time, int) {
	t := time.Date(2021, time.January, 6, 8, 0, 0, 0, time.Local)
	rand.Seed(time.Now().Unix())
	days := 365

	randomDay := rand.Intn(days)
	randomTime := t.Add(time.Hour * 24 * time.Duration(randomDay))

	randUserID := rand.Intn(10)
	return randomTime, randUserID
}

func CreateAndSaveJokes(joke model.Joke) error {

	filename := joke.ID + ".json"

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}

	dataBytes, err := json.MarshalIndent(joke, "", "   ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
func CreateAndSaveMessages(res string) (string, error) {

	filename := time.Now().Format("20060102150405") + ".csv"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return "", err
		}
	}

	// dataBytes, err := json.MarshalIndent(res, "", "   ")
	// if err != nil {
	// 	return "", err
	// }

	err = ioutil.WriteFile(filename, []byte(res), 0644)
	if err != nil {
		return "", err
	}

	return filename, nil
}
