package filestorage

import (
	"program/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var ms, _ = NewMongoStorage("mongodb://localhost:27017")

func TestFindID(t *testing.T) {

	var j = model.Joke{
		Title: "Buy cheese and bread for breakfast.",
		Body:  "and go away from me",
		Score: 1,
		ID:    "76h8ji",
	}

	_, err := ms.FindID("gesgsg1")
	assert.Equal(t, err, mongo.ErrNoDocuments)

	x, _ := ms.FindID("76h8ji")
	assert.Equal(t, j, x)
}

func TestFun(t *testing.T) {

	var j = "On the condition he gets to install windows.\n\n\n"

	r, _ := ms.Fun()
	r1 := r[0]
	assert.Equal(t, j, r1.Body)

}

func TestTextS(t *testing.T) {
	var s = "porcupinetree"

	_, err := ms.TextSearch(s)
	assert.Equal(t, err, nil)

}

func TestUpdateByID(t *testing.T) {

	var j = model.Joke{
		Body: "updaaat4e",
		ID:   "124124",
	}

	res, _ := ms.UpdateByID(j.Body, j.ID)

	assert.Equal(t, res.ModifiedCount, int64(1))

	var j2 = model.Joke{
		Body: "upd v.6",
		ID:   "124fagawg",
	}

	res1, _ := ms.UpdateByID(j2.Body, j2.ID)

	assert.Equal(t, res1.ModifiedCount, int64(0))

}
