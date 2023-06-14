package mongostorage_test

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"program/model"
	"program/storage/mongostorage"
	"program/testdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	file := "../../reddit_jokes.json"
	testdb.TestDB(file)
}

func initConnection() *mongostorage.MongoStorage {
	ms, _ := mongostorage.NewMongoStorage(os.Getenv("MONGODB_URI"))
	return ms
}

// /test
func TestFindID(t *testing.T) {
	ms := initConnection()
	j := model.Joke{
		Title: "What's the difference between a Jew in Nazi Germany and pizza ?",
		Body:  "Pizza doesn't scream when you put it in the oven .\n\nI'm so sorry.",
		Score: 0,
		ID:    "5tz4dd",
	}
	e := model.Joke{
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
		res, err := ms.FindID("5tz4dd")
		require.NoError(t, err)
		assert.Equal(t, j, res)
	})
}

func TestFun(t *testing.T) {
	ms := initConnection()
	var limit int64
	j := "Plagiarism. "

	res, _ := ms.Fun(limit)
	r := res[0]
	assert.Equal(t, j, r.Body)
}

func TestTextS(t *testing.T) {
	ms := initConnection()
	s := "porcupinetree"

	_, err := ms.TextSearch(s)
	require.Error(t, err)
	assert.Equal(t, err, mongo.ErrNoDocuments)
}

func TestUpdateByID(t *testing.T) {
	ms := initConnection()
	letters := "abcd"
	rand.Seed(time.Now().UnixNano())
	rand := strconv.Itoa(rand.Intn(1000) + 1)
	body := letters + string(rand)

	j := model.Joke{
		Body: body,
		ID:   "5tz1o1",
	}
	j2 := model.Joke{
		Body: "upd v.7",
		ID:   "124fagawg",
	}
	t.Run("IDexists", func(t *testing.T) {
		res, err := ms.UpdateByID(j.Body, j.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(1), res.ModifiedCount)
	})

	t.Run("NoID", func(t *testing.T) {
		res, err := ms.UpdateByID(j2.Body, j2.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), res.ModifiedCount)
	})
}

func TestLogin(t *testing.T) {
	ms := initConnection()
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
	ms := initConnection()
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
	ms := initConnection()
	var u model.User

	u.Username = "Denyska"
	err := ms.CreateUser(u)
	require.NoError(t, err)
}

func TestRandom(t *testing.T) {
	ms := initConnection()

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
	ms := initConnection()
	var u model.User
	u.Username = "D"
	token := "123"
	refreshToken := "1234"

	err := ms.UpdateTokens(token, refreshToken, u.Username)
	require.NoError(t, err)
}
