package mongostorage_test

import (
	"program/model"
	"program/storage/mongostorage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func Col() *mongostorage.MongoStorage {

	var ms, _ = mongostorage.NewMongoStorage("mongodb://localhost:27017")
	return ms
}
func TestFindID(t *testing.T) {
	ms := Col()
	var j = model.Joke{
		Title: "Breaking News: Bill Gates has agreed to pay for Trump's wall",
		Body:  "On the condition he gets to install windows.\n\n\n",
		Score: 48526,
		ID:    "5tn84z",
	}

	_, err := ms.FindID("gesgsg1")
	assert.Equal(t, err, mongo.ErrNoDocuments)

	x, _ := ms.FindID("5tn84z")
	assert.Equal(t, j, x)
}

func TestFun(t *testing.T) {
	ms := Col()
	var limit int64
	var j = "On the condition he gets to install windows.\n\n\n"

	r, _ := ms.Fun(limit)
	r1 := r[0]
	assert.Equal(t, j, r1.Body)

}

func TestTextS(t *testing.T) {
	ms := Col()
	var s = "porcupinetree"

	_, err := ms.TextSearch(s)
	require.Error(t, err)

}

func TestUpdateByID(t *testing.T) {
	ms := Col()
	var j = model.Joke{
		Body: "updaaat4e v.2",
		ID:   "5sfp26",
	}

	res, err := ms.UpdateByID(j.Body, j.ID)
	require.NoError(t, err)
	assert.Equal(t, res.ModifiedCount, int64(1))

	var j2 = model.Joke{
		Body: "upd v.7",
		ID:   "124fagawg",
	}

	res1, _ := ms.UpdateByID(j2.Body, j2.ID)

	assert.Equal(t, res1.ModifiedCount, int64(0))

}
