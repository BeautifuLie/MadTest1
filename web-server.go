package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"program/storage"
)

//type JokeSource interface {
//	LoadJokes() ([]Joke, map[string]Joke)
//	SaveJokes([]Joke, map[string]Joke)
//}
//
//var a JokeSource
//
//type Joke struct {
//	Title string `json:"title"`
//	Body  string `json:"body"`
//	Score int    `json:"score"`
//	ID    string `json:"id"`
//}
//
//type Server struct {
//	jokesStruct []Joke
//	jokesMap    map[string]Joke
//}

func main() {

	//s := Server{
	//	[]Joke{},
	//	map[string]Joke{},
	//}
	//
	//a = &s
	//jsonUnmarsh(&storage.S)

	storage.St = &storage.S
	myRouter := handleRequest(&storage.Server{})

	err := http.ListenAndServe(":9090", myRouter)
	if err != nil {
		panic(err)
	}

}

//func jsonUnmarsh(S *storage.Server) {
//
//	j, err := ioutil.ReadFile("reddit_jokes.json")
//	if err != nil {
//		fmt.Println("Error reading file", err)
//	}
//
//	err = json.Unmarshal(j, &storage.S.JokesStruct)
//	if err != nil {
//		fmt.Println("Error unmarshalling JSON", err)
//	}
//
//	sort.SliceStable(storage.S.JokesStruct, func(i, j int) bool {
//		return storage.S.JokesStruct[i].Score > storage.S.JokesStruct[j].Score
//	})
//
//	for _, j := range storage.S.JokesStruct {
//		storage.S.JokesMap[j.ID] = j
//
//	}
//
//}

//func (s *storage.Server) LoadJokes() ([]Joke, map[string]Joke) {
//	j, err := ioutil.ReadFile("reddit_jokes.json")
//	if err != nil {
//		fmt.Println("Error reading file", err)
//	}
//
//	err = json.Unmarshal(j, &s.jokesStruct)
//	if err != nil {
//		fmt.Println("Error unmarshalling JSON", err)
//	}
//
//	for _, j := range s.jokesStruct {
//		s.jokesMap[j.ID] = j
//	}
//	return s.jokesStruct, s.jokesMap
//}
//func (s *Server) Load(w http.ResponseWriter, r *http.Request) {
//	a.LoadJokes()
//	json.NewEncoder(w).Encode("File loaded")
//}
//
//func (s *Server) SaveJokes([]Joke, map[string]Joke) {
//
//	structBytes, err := json.MarshalIndent(s.jokesStruct, "", " ")
//	if err != nil {
//		fmt.Println("Error marshalling JSON", err)
//	}
//	err = ioutil.WriteFile("jokesStruct.json", structBytes, 0644)
//
//	mapBytes, err := json.MarshalIndent(s.jokesMap, "", " ")
//	if err != nil {
//		fmt.Println("Error marshalling JSON", err)
//	}
//	err = ioutil.WriteFile("jokesMap.json", mapBytes, 0644)
//}
//func (s *Server) Save(w http.ResponseWriter, r *http.Request) {
//	a.SaveJokes(s.jokesStruct, s.jokesMap)
//	if a.SaveJokes != nil {
//		err := json.NewEncoder(w).Encode("File saved")
//		if err != nil {
//			fmt.Println("Error saving file", err)
//		}
//	} else {
//		json.NewEncoder(w).Encode("Error saving file")
//	}
//}

func handleRequest(s *storage.Server) *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/jokes", homePage).Methods("GET")
	myRouter.HandleFunc("/jokes/method/save", Save)
	myRouter.HandleFunc("/jokes/method/load", Load)
	myRouter.HandleFunc("/jokes/funniest", getFunniestJokes)
	myRouter.HandleFunc("/jokes/random", getRandomJoke)
	myRouter.HandleFunc("/jokes", addJoke).Methods("POST")
	myRouter.HandleFunc("/jokes/{id}", getJokeByID)
	myRouter.HandleFunc("/jokes/search/{text}", getJokeByText)

	return myRouter
}

func getJokeByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	res, err := storage.ID(id, &storage.S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func getJokeByText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	text := vars["text"]
	res, err := storage.Text(text, &storage.S)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func getFunniestJokes(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}
	res, err1 := storage.Funniest(m, &storage.S)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func getRandomJoke(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err, "Error parsing query")
	}

	res, err1 := storage.Random(m, &storage.S)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func addJoke(w http.ResponseWriter, r *http.Request) {

	var j storage.Joke
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err1 := storage.Add(j, &storage.S)
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}

func Load(w http.ResponseWriter, r *http.Request) {

	_, err := storage.St.Load()
	if err != nil {
		http.Error(w, "Error loading file", 402)
	}

	json.NewEncoder(w).Encode("File loaded")

}

func Save(w http.ResponseWriter, r *http.Request) {

	err := storage.St.Save(storage.S.JokesStruct)
	if err != nil {
		http.Error(w, "oops", 400)
	}
	json.NewEncoder(w).Encode("File saved")

}
