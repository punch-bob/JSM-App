package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type JokeServer struct {
	store *JokesStore
}

func NewJokeServer() *JokeServer {
	store := NewJokesStore()
	return &JokeServer{store: store}
}

func (server *JokeServer) JokeListHandler(write http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(server.store.GetAllJokes())
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) UpdateJokeRatingHandler(write http.ResponseWriter, request *http.Request) {
	type RequestReact struct {
		Id       int    `json:"id"`
		Reaction string `json:"reaction"`
	}

	dec := json.NewDecoder(request.Body)
	var rr RequestReact
	if err := dec.Decode(&rr); err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	var jokeRating int
	if rr.Reaction == "increase" {
		jokeRating = server.store.IncreaseRating(rr.Id)
	} else if rr.Reaction == "decrease" {
		jokeRating = server.store.DecreaseRating(rr.Id)
	}

	js, err := json.Marshal(jokeRating)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) Handler(write http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		path := strings.Trim(request.URL.Path, "/")
		splitedPath := strings.Split(path, "/")
		if len(splitedPath) != 2 {
			http.Error(write, "expect /hello/<your name> in task handler", http.StatusBadRequest)
			return
		}

		username := strings.ReplaceAll(splitedPath[1], "_", " ")
		helloMessage := fmt.Sprintf("Hello %s!", username)
		js, err := json.Marshal(helloMessage)
		if err != nil {
			http.Error(write, err.Error(), http.StatusInternalServerError)
			return
		}

		write.Header().Set("Access-Control-Allow-Origin", "*")
		write.Header().Set("Content-Type", "application/json")
		write.Write(js)
	} else {
		http.Error(write, fmt.Sprintf("expect method GET at /hello/<your name>/, got %v", request.Method), http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	var joke1 Joke
	joke1.Id = 0
	joke1.Text = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	joke1.Rating = 10
	joke1.Tags = make([]string, 2)
	joke1.Tags[0] = "test"
	joke1.Tags[1] = "test1"
	joke1.AuthorName = "jsm"
	joke1.Date = "2022.07.12 18:50:23"

	var joke2 Joke
	joke2.Id = 1
	joke2.Text = "abbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\naaaaaa\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	joke2.Rating = -10
	joke2.Tags = make([]string, 3)
	joke2.Tags[0] = "test2"
	joke2.Tags[1] = "test3"
	joke2.Tags[2] = "test4"
	joke2.AuthorName = "jsm"
	joke2.Date = "2022.07.12 18:50:23"

	router := mux.NewRouter()
	server := NewJokeServer()

	server.store.Store[0] = joke1
	server.store.Store[1] = joke2

	router.HandleFunc("/joke_list/", server.JokeListHandler).Methods("GET")
	router.HandleFunc("/update_rating/", server.UpdateJokeRatingHandler).Methods("POST")

	//Header sets
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "HEAD", "PUT", "OPTIONS"})

	http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(router))
}
