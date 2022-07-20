package main

import (
	"sync"
	"time"
)

type Joke struct {
	Id         int           `json:"id"`
	Text       string        `json:"text"`
	Rating     int           `json:"rate"`
	Tags       []string      `json:"tags"`
	AuthorName string        `json:"author_name"`
	Date       time.Time     `json:"date"`
	WhoRated   map[int]uint8 `json:"who_rated"`
}

type JokesStore struct {
	Store map[int]Joke
	sync.Mutex
	CurId           int
	GeneratedJokeId int
}

func NewJokesStore() *JokesStore {
	jokesStore := &JokesStore{}
	jokesStore.CurId = 0
	jokesStore.Store = make(map[int]Joke)
	jokesStore.GeneratedJokeId = -1
	return jokesStore
}

func (js *JokesStore) CreateJoke(text string, tags []string, authorName string) Joke {
	js.Lock()
	defer js.Unlock()

	joke := Joke{
		Id:         js.CurId,
		Text:       text,
		Rating:     0,
		AuthorName: authorName,
		Date:       time.Now(),
		WhoRated:   make(map[int]uint8)}
	joke.Tags = make([]string, len(tags))
	copy(joke.Tags, tags)

	js.Store[js.CurId] = joke
	js.CurId++
	return joke
}

func (js *JokesStore) IncreaseRating(uid, id int) int {
	js.Lock()
	defer js.Unlock()

	joke := js.Store[id]
	if val, ok := joke.WhoRated[uid]; ok {
		if val == 2 {
			return joke.Rating
		}
	} else {
		joke.WhoRated[uid] = 1
	}

	joke.Rating++
	joke.WhoRated[uid]++
	js.Store[id] = joke
	return joke.Rating
}

func (js *JokesStore) DecreaseRating(uid, id int) int {
	js.Lock()
	defer js.Unlock()

	joke := js.Store[id]
	if val, ok := joke.WhoRated[uid]; ok {
		if val == 0 {
			return joke.Rating
		}
	} else {
		joke.WhoRated[uid] = 1
	}

	joke.Rating--
	joke.WhoRated[uid]--
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
	js.Lock()
	defer js.Unlock()

	jokes := make([]Joke, 0, len(js.Store))

	for id, joke := range js.Store {
		if id != js.GeneratedJokeId {
			jokes = append(jokes, joke)
		}
	}
	return jokes
}

func (js *JokesStore) GetDailyJoke() Joke {
	js.Lock()
	defer js.Unlock()

	var dailyJoke Joke

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	const MaxInt = int(^uint(0) >> 1)
	maxRate := -MaxInt - 1

	for _, joke := range js.Store {
		jokeDate := joke.Date
		dif := today.Sub(jokeDate)
		if joke.Rating > maxRate && 0 < dif && dif < 24*time.Hour {
			dailyJoke = joke
			maxRate = joke.Rating
		}
	}
	return dailyJoke
}

func (js *JokesStore) GetGeneratedJoke() Joke {
	return js.Store[js.GeneratedJokeId]
}
