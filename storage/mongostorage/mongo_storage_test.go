package mongostorage_test

import (
	"errors"
	"program/model"
	"program/storage/mongostorage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewMongoStorage(t *testing.T) {

	var _, err = mongostorage.NewMongoStorage("mongodb://localhost:27018")
	require.Error(t, err)

}
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
	var e = model.Joke{
		Title: "",
		Body:  "",
		Score: 0,
		ID:    "",
	}
	t.Run("IDnotExists", func(t *testing.T) {

		res, err := ms.FindID("gesgsg1")
		require.Error(t, err)
		assert.Equal(t, err, mongo.ErrNoDocuments)
		assert.Equal(t, e, res)
	})
	t.Run("IDexists", func(t *testing.T) {

		res, err := ms.FindID("5tn84z")
		require.NoError(t, err)
		assert.Equal(t, j, res)
	})

}

func TestFun(t *testing.T) {
	ms := Col()
	var limit int64
	var j = "On the condition he gets to install windows.\n\n\n"

	res, _ := ms.Fun(limit)
	r := res[0]
	assert.Equal(t, j, r.Body)

}

func TestTextS(t *testing.T) {
	ms := Col()
	var s = "porcupinetree"

	_, err := ms.TextSearch(s)
	require.Error(t, err)
	assert.Equal(t, err, mongo.ErrNoDocuments)

}

func TestUpdateByID(t *testing.T) {
	ms := Col()
	var j = model.Joke{
		Body: "updaaat4e v.2",
		ID:   "1234",
	}
	var j2 = model.Joke{
		Body: "upd v.7",
		ID:   "124fagawg",
	}
	t.Run("IDexists", func(t *testing.T) {

		res, err := ms.UpdateByID(j.Body, j.ID)
		require.NoError(t, err)
		assert.NotEqual(t, res.ModifiedCount, int64(1))
	})

	t.Run("NoID", func(t *testing.T) {

		res, err := ms.UpdateByID(j2.Body, j2.ID)
		assert.NoError(t, err)
		assert.Equal(t, res.ModifiedCount, int64(0))
	})

}

func TestLogin(t *testing.T) {
	ms := Col()
	var u model.User
	t.Run("NoUser", func(t *testing.T) {
		u.Username = "zxc"
		res, err := ms.LoginUser(u)
		require.Error(t, err)
		assert.NotEqual(t, u.Username, res.Username)
	})
	t.Run("UserExists", func(t *testing.T) {
		u.Username = "Denys"
		res, err := ms.LoginUser(u)
		require.NoError(t, err)
		assert.Equal(t, u.Username, res.Username)
	})

}
func TestIsExists(t *testing.T) {
	ms := Col()
	var u model.User
	noExist := errors.New("this username already exists")
	t.Run("NoUser", func(t *testing.T) {
		u.Username = "zxc"
		err := ms.IsExists(u)
		require.NoError(t, err)

	})
	t.Run("UserExists", func(t *testing.T) {
		u.Username = "Denys"
		err := ms.IsExists(u)
		require.Error(t, err)
		assert.Equal(t, err, noExist)

	})

}
func TestCreateUser(t *testing.T) {
	ms := Col()
	var u model.User

	u.Username = "Denyska"
	err := ms.CreateUser(u)
	require.NoError(t, err)

}
func TestRandom(t *testing.T) {
	ms := Col()

	t.Run("RandomPositiveLimit", func(t *testing.T) {
		limit := 5
		res, err := ms.Random(limit)
		require.NoError(t, err)
		assert.Equal(t, limit, len(res))

	})
	t.Run("RandomNegativeLimit", func(t *testing.T) {
		limit := -5
		res, _ := ms.Random(limit)
		assert.Equal(t, 0, len(res))

	})

}
func TestUpdateTOkens(t *testing.T) {
	ms := Col()
	var u model.User
	u.Username = "D"
	token := "123"
	refreshToken := "1234"

	err := ms.UpdateTokens(token, refreshToken, u.Username)
	require.NoError(t, err)

}
