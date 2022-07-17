package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

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

func (server *JokeServer) jokeListHandler(write http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(server.store.GetAllJokes())
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) updateJokeRatingHandler(write http.ResponseWriter, request *http.Request) {
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

func (server *JokeServer) dailyJokeHandler(write http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(server.store.GetDailyJoke())
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) generatedJokeHandler(write http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(server.store.GetGeneratedJoke())
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) generateJoke() int {
	rand.Seed(time.Now().UnixNano())

	requestBody, err := json.Marshal(map[string]string{
		"text": "поручик Ржевский пришел к Штирлицу",
	})
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("POST", "https://api.aicloud.sbercloud.ru/public/v1/public_inference/gpt3/predict", bytes.NewBuffer(requestBody))
	request.Header = http.Header{
		"Connection":      {"keep-alive"},
		"Accept":          {"application/json"},
		"Server":          {"istio-envoy"},
		"Accept-encoding": {"gzip, deflate, br"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 OPR/82.0.4227.43"},
		"Content-Type":    {"application/json"},
		"Origin":          {"https://russiannlp.github.io"},
		"Referer":         {"https://russiannlp.github.io/"},
		"Sec-Fetch-Site":  {"same-site"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Dest":  {"empty"},
		"Accept-Language": {"en-US,en;q=0.9,es-AR;q=0.8,es;q=0.7"},
	}
	if err != nil {
		log.Fatalln(err)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	type AIResponse struct {
		Predictions string `json:"predictions"`
	}

	dec := json.NewDecoder(response.Body)
	var airesp AIResponse
	if err := dec.Decode(&airesp); err != nil {
		log.Fatalln(err)
	}

	jokeText := airesp.Predictions
	tags := []string{"AI"}
	genJokeId := server.store.CreateJoke(string(jokeText), tags, "ruGPT-3 XL")

	return genJokeId
}

func (server *JokeServer) createJokeHandler(write http.ResponseWriter, request *http.Request) {
	type JokeLayout struct {
		Text       string   `json:"text"`
		Tags       []string `json:"tags"`
		AuthorName string   `json:"author_name"`
	}

	dec := json.NewDecoder(request.Body)
	var jl JokeLayout
	err := dec.Decode(&jl)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	server.store.CreateJoke(jl.Text, jl.Tags, jl.AuthorName)
}

func (server *JokeServer) findJokeByTagsHendler(write http.ResponseWriter, request *http.Request) {
	type RequestTemplate struct {
		Tags []string `json:"tags"`
	}

	dec := json.NewDecoder(request.Body)
	var rt RequestTemplate
	err := dec.Decode(&rt)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	jokeList := server.store.GetJokesByTags(rt.Tags)
	js, err := json.Marshal(jokeList)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func main() {
	var joke1 Joke
	joke1.Id = 0
	joke1.Text = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	joke1.Rating = 10
	joke1.Tags = make([]string, 2)
	joke1.Tags[0] = "test"
	joke1.Tags[1] = "test1"
	joke1.AuthorName = "champ"
	joke1.Date = time.Now().Add(-24 * time.Hour)

	var joke2 Joke
	joke2.Id = 1
	joke2.Text = "тестbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\naaaaaa\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	joke2.Rating = -10
	joke2.Tags = make([]string, 3)
	joke2.Tags[0] = "test2"
	joke2.Tags[1] = "test3"
	joke2.Tags[2] = "test4"
	joke2.AuthorName = "loser"
	joke2.Date = time.Now().Add(-24 * time.Hour)

	router := mux.NewRouter()
	server := NewJokeServer()

	server.store.Store[0] = joke1
	server.store.Store[1] = joke2
	server.store.CurId = 2
	server.store.GeneratedJokeId = server.generateJoke()

	router.HandleFunc("/joke_list/", server.jokeListHandler).Methods("GET")
	router.HandleFunc("/update_rating/", server.updateJokeRatingHandler).Methods("POST")
	router.HandleFunc("/daily_joke/", server.dailyJokeHandler).Methods("GET")
	router.HandleFunc("/generated_joke/", server.generatedJokeHandler).Methods("GET")
	router.HandleFunc("/create_joke/", server.createJokeHandler).Methods("POST")
	router.HandleFunc("/joke_by_tags/", server.findJokeByTagsHendler).Methods("POST")

	//Header sets
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "HEAD", "PUT", "OPTIONS"})

	http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(router))
}
