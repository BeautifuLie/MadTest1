package joker_test

import (
	"net/url"
	"program/joker"
	"program/model"
	"program/storage"
	mocks "program/storage/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	j := model.Joke{
		ID:   "test",
		Body: "haha",
	}

	store.EXPECT().FindID("test").Return(j, nil)

	s := joker.NewServer(store)
	got, err := s.ID("test")

	require.NoError(t, err)
	require.Equal(t, j, got)

}

func TestFunniest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	j1 := model.Joke{
		Title: "fawfaw",
		Body:  "haha",
		ID:    "1q2w3e",
		Score: 1,
	}
	j2 := model.Joke{
		Title: "other",
		Body:  "haha1",
		ID:    "4r5t6y",
		Score: 2,
	}
	j := []model.Joke{j1, j2}

	store.EXPECT().Fun().Return(j, nil)

	s := joker.NewServer(store)

	m, _ := url.ParseQuery("limit=2")
	got, err := s.Funniest(m)
	require.NoError(t, err)
	require.Equal(t, j, got)
	j3 := j[0]
	assert.Equal(t, "haha", j3.Body)

}

func TestRandom(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	j1 := model.Joke{Title: "fawfaw", Body: "haha", ID: "1q2w3e", Score: 1}
	j2 := model.Joke{Title: "other", Body: "haha1", ID: "4r5t6y", Score: 2}
	j3 := model.Joke{Title: "hhzrh", Body: "4gzgz", ID: "g8g9j", Score: 5}
	j4 := model.Joke{Title: "ogrgzr", Body: "hz4hz", ID: "0g7g8f", Score: 10}
	j := []model.Joke{j1, j2, j3, j4}

	store.EXPECT().Random().Times(2).Return(j, nil)

	s := joker.NewServer(store)

	m, _ := url.ParseQuery("limit=4")
	r1, err := s.Random(m)

	require.NoError(t, err)

	r2, _ := s.Random(m)
	require.NotEqual(t, r1, r2)

}
func TestText(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	// j1 := model.Joke{Title: "fawfaw", Body: "haha", ID: "1q2w3e", Score: 1}
	// j2 := model.Joke{Title: "other", Body: "haha1", ID: "4r5t6y", Score: 2}

	j := []model.Joke{}

	store.EXPECT().TextSearch("jira").Return(j, nil).Times(1)

	s := joker.NewServer(store)
	got, err := s.Text("jira")

	require.NoError(t, err)
	require.Equal(t, j, got)

}

func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	j := model.Joke{Title: "fawfaw", Body: "haha", ID: "1q2w3e", Score: 1}

	store.EXPECT().Save(j).Return(nil)
	s := joker.NewServer(store)
	_, err := s.Add(j)
	require.NoError(t, err)

}

func TestUpdateB(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	j := model.Joke{Title: "fawfaw", Body: "tratata", ID: "1q2w3e", Score: 1}
	a := &mongo.UpdateResult{}
	store.EXPECT().UpdateByID(j.Body, j.ID).Return(a, storage.ErrNoJokes)
	s := joker.NewServer(store)
	_, err := s.Update(j, j.ID)
	require.Error(t, err)

}
