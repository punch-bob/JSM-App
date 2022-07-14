package main

import (
	"sync"
	"time"
)

type Joke struct {
	Id         int      `json:"id"`
	Text       string   `json:"text"`
	Rating     int      `json:"rate"`
	Tags       []string `json:"tags"`
	AuthorName string   `json:"author_name"`
	Date       string   `json:"date"`
}

type JokesStore struct {
	Store map[int]Joke
	sync.Mutex
	CurId int
}

func NewJokesStore() *JokesStore {
	jokesStore := &JokesStore{}
	jokesStore.CurId = 0
	jokesStore.Store = make(map[int]Joke)
	return jokesStore
}

func (js *JokesStore) CreateJoke(text string, tags []string, authorName string) {
	js.Lock()
	defer js.Unlock()

	joke := Joke{
		Id:         js.CurId,
		Text:       text,
		Rating:     0,
		AuthorName: authorName,
		Date:       time.Now().Format("2022.07.12 15:04:05")}
	joke.Tags = make([]string, len(tags))
	copy(joke.Tags, tags)

	js.Store[js.CurId] = joke
	js.CurId++
}

func (js *JokesStore) IncreaseRating(id int) int {
	js.Lock()
	defer js.Unlock()

	joke := js.Store[id]
	joke.Rating++
	js.Store[id] = joke

	return joke.Rating
}

func (js *JokesStore) DecreaseRating(id int) int {
	js.Lock()
	defer js.Unlock()

	joke := js.Store[id]
	joke.Rating--
	js.Store[id] = joke

	return joke.Rating
}

func (js *JokesStore) GetJokesByTags(tags []string) []Joke {
	js.Lock()
	defer js.Unlock()

	var jokes []Joke

jokeloop:
	for _, joke := range js.Store {
		for _, jokeTag := range joke.Tags {
			for _, tag := range tags {
				if tag == jokeTag {
					jokes = append(jokes, joke)
					continue jokeloop
				}
			}
		}
	}
	return jokes
}

func (js *JokesStore) GetAllJokes() []Joke {
	jokes := make([]Joke, 0, len(js.Store))

	for _, joke := range js.Store {
		jokes = append(jokes, joke)
	}
	return jokes
}
