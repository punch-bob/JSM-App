package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var lastUpdate time.Time
var SERVER_PORT = os.Getenv("SERVER_PORT")

type JokeServer struct {
	store *JokesStore
}

var jokeThemes = []string{"1 апреля", "14 февраля", "23 февраля", "8 марта", "Абрамович",
	"Адам и Ева", "армянское радио", "Баба-Яга", "Березовский", "Билл Гейтс", "блондинки", "богатыри", "Брежнев",
	"британские ученые", "Буратино", "Валуев", "веган", "Винни-Пух", "вовочка",
	"немец, американец и русский", "гаи", "геи", "девушки", "Дед Мороз", "Донцова",
	"золотая рыбка", "Иван-царевич", "Искусственный интеллект", "Каренина",
	"Карлсон", "Колобок", "Красная Шапочка", "чиновник", "Илон Маск", "Кот Шрёдингера",
	"Куклачев", "милиция", "Мавроди", "муж и жена", "Навальный",
	"наркоман Павлик", "новые русские", "Обама", "Перельман", "случай в поезде",
	"программист", "Прохоров", "Пушкин", "Рабинович", "поручик Ржевский",
	"сантехник", "сбербанк", "Сталин",
	"студент и профессор", "Сусанин", "тёща", "Трамп", "Чак Норрис",
	"Чапаев", "Чебурашка и корокодил Гена", "чукча", "Шерлок Холмс", "Штирлиц"}

var jokeAction = []string{"шел как-то по лесу", "встретил НЛО", "оказался в тылу", "зашел в бар", "выстрелил в ногу",
	"сел в машину", "прибежал домой и  увидел там", "забежал в ванную", "попал в плен", "уходил от погони", "играли в нарды",
	"зашли как-то в лифт", "выпили по стакану пива", "неожиданно в баре материализуется",
	"в белом плаще с кровавым подбоем, шаркающей кавалерийской походкой заходит в бар", "собрались на рыбалку", "ограбили"}

func NewJokeServer() *JokeServer {
	store := newJokesStore()
	return &JokeServer{
		store: store,
	}
}

func (server *JokeServer) jokeListHandler(write http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(server.store.getJokeList())
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
		Uid      int    `json:"uid"`
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
		jokeRating = server.store.increaseRating(rr.Uid, rr.Id)
	} else if rr.Reaction == "decrease" {
		jokeRating = server.store.decreaseRating(rr.Uid, rr.Id)
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
	js, err := json.Marshal(server.store.getDailyJoke())
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) generatedJokeHandler(write http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(server.generateJoke())
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) generateJoke() Joke {
	if lastUpdate.IsZero() || time.Since(lastUpdate) >= time.Hour*24 {
		rand.Seed(time.Now().UnixNano())

		reqText := jokeThemes[rand.Intn(len(jokeThemes))] + " " + jokeAction[rand.Intn(len(jokeAction))]
		requestBody, err := json.Marshal(map[string]string{
			"text": reqText,
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
		genJoke := server.store.createJoke(string(jokeText), tags, "ruGPT-3 XL", -1)

		lastUpdate = time.Now()
		lastUpdate = time.Date(lastUpdate.Year(), lastUpdate.Month(), lastUpdate.Day(), 0, 0, 0, 0, lastUpdate.Location())
		server.store.GeneratedJoke = genJoke
		return genJoke
	} else {
		return server.store.GeneratedJoke
	}
}

func (server *JokeServer) createJokeHandler(write http.ResponseWriter, request *http.Request) {
	type JokeLayout struct {
		Text       string   `json:"text"`
		Tags       []string `json:"tags"`
		AuthorName string   `json:"author_name"`
		UserId     int      `json:"uid"`
	}

	dec := json.NewDecoder(request.Body)
	var jl JokeLayout
	err := dec.Decode(&jl)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(server.store.createJoke(jl.Text, jl.Tags, jl.AuthorName, jl.UserId))
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) findJokeByTagsHandler(write http.ResponseWriter, request *http.Request) {
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

	jokeList := server.store.getJokesByTagList(rt.Tags)
	js, err := json.Marshal(jokeList)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *JokeServer) deleteJokeHandler(write http.ResponseWriter, request *http.Request) {
	type RequestTemplate struct {
		Id int `json:"id"`
	}

	dec := json.NewDecoder(request.Body)
	var rt RequestTemplate
	err := dec.Decode(&rt)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	server.store.deleteJoke(rt.Id)
}

func (server *JokeServer) findJokeByUIDHandler(write http.ResponseWriter, request *http.Request) {
	type RequestTemplate struct {
		Uid int `json:"uid"`
	}

	dec := json.NewDecoder(request.Body)
	var rt RequestTemplate
	err := dec.Decode(&rt)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	jokeList := server.store.getJokesByUID(rt.Uid)
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
	log.Println("start")

	router := mux.NewRouter()
	server := NewJokeServer()
	defer server.store.db.Close()

	router.HandleFunc("/joke_list/", server.jokeListHandler).Methods("GET")
	router.HandleFunc("/update_rating/", server.updateJokeRatingHandler).Methods("POST")
	router.HandleFunc("/daily_joke/", server.dailyJokeHandler).Methods("GET")
	router.HandleFunc("/generated_joke/", server.generatedJokeHandler).Methods("GET")
	router.HandleFunc("/create_joke/", server.createJokeHandler).Methods("POST")
	router.HandleFunc("/joke_by_tags/", server.findJokeByTagsHandler).Methods("POST")
	router.HandleFunc("/joke_by_uid/", server.findJokeByUIDHandler).Methods("POST")
	router.HandleFunc("/delete_joke/", server.deleteJokeHandler).Methods("DELETE")

	//Header sets
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "HEAD", "PUT", "OPTIONS"})

	http.ListenAndServe(":"+SERVER_PORT, handlers.CORS(headersOk, originsOk, methodsOk)(router))
}
