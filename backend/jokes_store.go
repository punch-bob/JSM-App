package main

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

type Joke struct {
	Id         int       `json:"id"`
	Text       string    `json:"text"`
	Rating     int       `json:"rate"`
	Tags       []string  `json:"tags"`
	AuthorName string    `json:"author_name"`
	Date       time.Time `json:"date"`
}

type JokesStore struct {
	Store map[int]Joke
	sync.Mutex
	GeneratedJokeId int
	db              *sql.DB
}

func NewJokesStore() *JokesStore {
	Store := make(map[int]Joke)
	GeneratedJokeId := -1
	// dataSourceName := os.Getenv("DB_LOGIN") + ":" + os.Getenv("DB_PASSWORD") + "@/" + os.Getenv("DB_NAME")
	dataSourceName := "root:password@/joke_db"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	LoadJokesFromDB(Store, db)
	return &JokesStore{
		Store:           Store,
		GeneratedJokeId: GeneratedJokeId,
		db:              db,
	}
}

func LoadJokeTagsFromDB(id int, db *sql.DB) []string {
	var tag []string
	// tags, err := db.Query("select tag from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME") + " where joke_id = ?", id)
	tags, err := db.Query("select tag from joke_db.tag where joke_id = ?", id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer tags.Close()

	for tags.Next() {
		var newTag string
		err := tags.Scan(&newTag)
		if err != nil {
			log.Println(err)
			continue
		}
		tag = append(tag, newTag)
	}
	return tag
}

func LoadJokesFromDB(jokeStore map[int]Joke, db *sql.DB) {
	// rows, err := db.Query("select * from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME"))
	rows, err := db.Query("select * from joke_db.joke")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		joke := Joke{}
		var date string
		err := rows.Scan(&joke.Id, &joke.Text, &joke.Rating, &joke.AuthorName, &date)
		if err != nil {
			log.Println(err)
			continue
		}
		joke.Date, err = time.Parse(time.RFC3339, date)
		if err != nil {
			log.Println(err)
		}
		joke.Tags = LoadJokeTagsFromDB(joke.Id, db)
		jokeStore[joke.Id] = joke
	}

}

func (js *JokesStore) UpdateRatingInBD(id, rating int) {
	// updateRating, err := js.db.Prepare("update "+os.Getenv("TABLE_NAME")+" set rating = ? where id = ?")
	updateRating, err := js.db.Prepare("update joke set rating = ? where id = ?")
	if err != nil {
		log.Println(err)
	}
	_, _ = updateRating.Exec(rating, id)
}

func (js *JokesStore) UpdateReactInBD(uid, id, react int) {
	// updateRating, err := js.db.Prepare("update "+os.Getenv("TABLE_NAME")+" set react = ? where joke_id = ? and user_id = ?")
	updateReact, err := js.db.Prepare("update joke_rating set react = ? where joke_id = ? and user_id = ?")
	if err != nil {
		log.Println(err)
	}
	_, _ = updateReact.Exec(react, id, uid)
}

func (js *JokesStore) AddJokeInDB(joke Joke) int {
	var (
		text       = joke.Text
		rating     = joke.Rating
		authorName = joke.AuthorName
		date       = joke.Date.Format(time.RFC3339)
	)
	// res, err := js.db.Exec("insert into "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" (text, rating, author_name, date) values (?, ?, ?, ?)", text, rating, authorName, date)
	res, err := js.db.Exec("insert into joke_db.joke (text, rating, author_name, date) values (?, ?, ?, ?)", text, rating, authorName, date)
	if err != nil {
		log.Println(err)
		return -1
	}
	id, _ := res.LastInsertId()
	joke.Id = int(id)

	for _, tag := range joke.Tags {
		// _, err := js.db.Exec("insert into "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" l(joke_id, tag) values (?, ?)", id, tag)
		_, err := js.db.Exec("insert into joke_db.tag (joke_id, tag) values (?, ?)", id, tag)
		if err != nil {
			log.Println(err)
		}
	}
	return int(id)
}

func (js *JokesStore) CreateJoke(text string, tags []string, authorName string) int {
	js.Lock()
	defer js.Unlock()

	joke := Joke{
		Id:         0,
		Text:       text,
		Rating:     0,
		AuthorName: authorName,
		Date:       time.Now()}
	joke.Tags = make([]string, len(tags))
	copy(joke.Tags, tags)

	js.AddJokeInDB(joke)

	js.Store[joke.Id] = joke
	return joke.Id
}

func (js *JokesStore) IncreaseRating(uid, id int) int {
	js.Lock()
	defer js.Unlock()

	joke := js.Store[id]
	// stmt, err := js.db.Prepare("select react from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" where joke_id = ? AND user_id = ?")
	stmt, err := js.db.Prepare("select react from joke_bd.joke_rating where joke_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
	}
	react := stmt.QueryRow(id, uid)
	var value int
	err = react.Scan(&value)
	if err == sql.ErrNoRows {
		// js.db.Exec("insert into "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" (joke_id, user_id, react) values (?, ?, ?)", id, uid, 1)
		js.db.Exec("insert into joke_db.tag (joke_id, user_id, react) values (?, ?, ?)", id, uid, 1)
		value = 1
	}

	if value == 2 {
		return joke.Rating
	}

	joke.Rating++
	js.UpdateRatingInBD(id, joke.Rating)
	js.UpdateReactInBD(uid, id, value+1)
	js.Store[id] = joke
	return joke.Rating
}

func (js *JokesStore) DecreaseRating(uid, id int) int {
	js.Lock()
	defer js.Unlock()

	joke := js.Store[id]
	// stmt, err := js.db.Prepare("select react from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" where joke_id = ? AND user_id = ?")
	stmt, err := js.db.Prepare("select react from joke_bd.joke_rating where joke_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
	}
	react := stmt.QueryRow(id, uid)
	var value int
	err = react.Scan(&value)
	if err == sql.ErrNoRows {
		// js.db.Exec("insert into "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" (joke_id, user_id, react) values (?, ?, ?)", id, uid, 1)
		js.db.Exec("insert into joke_db.tag (joke_id, user_id, react) values (?, ?, ?)", id, uid, 1)
		value = 1
	}

	if value == 0 {
		return joke.Rating
	}

	joke.Rating--
	js.UpdateRatingInBD(id, joke.Rating)
	js.UpdateReactInBD(uid, id, value-1)
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
