package main

import (
	"database/sql"
	"log"
	"strconv"
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
	sync.Mutex
	GeneratedJokeId int
	db              *sql.DB
}

func NewJokesStore() *JokesStore {
	GeneratedJokeId := -1
	// dataSourceName := os.Getenv("DB_LOGIN") + ":" + os.Getenv("DB_PASSWORD") + "@/" + os.Getenv("DB_NAME")
	dataSourceName := "root:password@/joke_db"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	return &JokesStore{
		GeneratedJokeId: GeneratedJokeId,
		db:              db,
	}
}

func removeDuplicateValues(intSlice []int) []int {
	keys := make(map[int]bool)
	var list []int

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

func GetJokeFromDB(db *sql.DB, id int) Joke {
	// row := db.QueryRow("select * from " + os.Getenv("DB_NAME") + "."+os.Getenv("TABLE_NAME"))
	row := db.QueryRow("select * from joke_db.joke where id = ?", id)
	joke := Joke{}
	err := row.Scan(&joke.Id, &joke.Text, &joke.Rating, &joke.AuthorName, &joke.Date)
	if err != nil {
		log.Println(err)
	}
	joke.Tags = LoadJokeTagsFromDB(id, db)
	return joke
}

func ConvertRowToJoke(db *sql.DB, row *sql.Rows) Joke {
	joke := Joke{}
	err := row.Scan(&joke.Id, &joke.Text, &joke.Rating, &joke.AuthorName, &joke.Date)
	if err != nil {
		log.Println(err)
	}
	joke.Tags = LoadJokeTagsFromDB(joke.Id, db)
	return joke
}

func LoadJokesFromDB(db *sql.DB) []Joke {
	var jokeStore []Joke
	// rows, err := db.Query("select * from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME"))
	rows, err := db.Query("select * from joke_db.joke")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		joke := Joke{}
		err := rows.Scan(&joke.Id, &joke.Text, &joke.Rating, &joke.AuthorName, &joke.Date)
		if err != nil {
			continue
		}
		joke.Tags = LoadJokeTagsFromDB(joke.Id, db)
		jokeStore = append(jokeStore, joke)
	}
	return jokeStore
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
		date       = joke.Date
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
	joke := Joke{
		Id:         0,
		Text:       text,
		Rating:     0,
		AuthorName: authorName,
		Date:       time.Now().Format("2006-01-02")}
	joke.Tags = make([]string, len(tags))
	copy(joke.Tags, tags)

	id := js.AddJokeInDB(joke)

	return id
}

func (js *JokesStore) IncreaseRating(uid, id int) int {
	joke := GetJokeFromDB(js.db, id)

	// stmt, err := js.db.Prepare("select react from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" where joke_id = ? AND user_id = ?")
	stmt, err := js.db.Prepare("select react from joke_db.joke_rating where joke_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
	}
	react := stmt.QueryRow(id, uid)
	value := 1
	err = react.Scan(&value)
	if err == sql.ErrNoRows {
		// js.db.Exec("insert into "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" (joke_id, user_id, react) values (?, ?, ?)", id, uid, 1)
		js.db.Exec("insert into joke_db.joke_rating (joke_id, user_id, react) values (?, ?, ?)", id, uid, value)
	}

	if value == 2 {
		return joke.Rating
	}

	joke.Rating++
	js.UpdateRatingInBD(id, joke.Rating)
	js.UpdateReactInBD(uid, id, value+1)
	// js.Store[id] = joke
	return joke.Rating
}

func (js *JokesStore) DecreaseRating(uid, id int) int {
	js.Lock()
	defer js.Unlock()

	// joke := js.Store[id]
	joke := GetJokeFromDB(js.db, id)

	// stmt, err := js.db.Prepare("select react from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" where joke_id = ? AND user_id = ?")
	stmt, err := js.db.Prepare("select react from joke_db.joke_rating where joke_id = ? and user_id = ?")
	if err != nil {
		log.Println(err)
	}
	react := stmt.QueryRow(id, uid)
	value := 1
	err = react.Scan(&value)
	if err == sql.ErrNoRows {
		// js.db.Exec("insert into "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" (joke_id, user_id, react) values (?, ?, ?)", id, uid, 1)
		js.db.Exec("insert into joke_db.joke_rating (joke_id, user_id, react) values (?, ?, ?)", id, uid, value)
	}

	if value == 0 {
		return joke.Rating
	}

	joke.Rating--
	js.UpdateRatingInBD(id, joke.Rating)
	js.UpdateReactInBD(uid, id, value-1)
	// js.Store[id] = joke
	return joke.Rating
}

func GetJokeIdByTag(db *sql.DB, tag string) []int {
	var jokeId []int

	// jokes, err := js.db.Query("select joke_id from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+"where tag = ?", tag)
	rows, err := db.Query("select joke_id from joke_db.tag where tag = ?", tag)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var value int
		err = rows.Scan(&value)
		if err != nil {
			continue
		}

		jokeId = append(jokeId, value)
	}
	return jokeId
}

func (js *JokesStore) GetJokesByTags(tags []string) []Joke {
	js.Lock()
	defer js.Unlock()

	var jokes []Joke
	var jokeId []int

	for _, tag := range tags {
		jokeId = append(jokeId, GetJokeIdByTag(js.db, tag)...)
		jokeId = removeDuplicateValues(jokeId)
	}

	for _, id := range jokeId {
		jokes = append(jokes, GetJokeFromDB(js.db, id))
	}
	return jokes
}

func (js *JokesStore) GetAllJokes() []Joke {
	js.Lock()
	defer js.Unlock()

	jokes := LoadJokesFromDB(js.db)
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

	// jokes, err := js.db.Query("select * from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME"))
	jokes, err := js.db.Query("select * from joke_db.joke")
	if err != nil {
		log.Println(err)
	}
	defer jokes.Close()

	for jokes.Next() {
		joke := ConvertRowToJoke(js.db, jokes)
		jokeDate := joke.Date
		year, _ := strconv.Atoi(jokeDate[:len(jokeDate)-6])
		month, _ := strconv.Atoi(jokeDate[len(jokeDate)-5 : len(jokeDate)-3])
		day, _ := strconv.Atoi(jokeDate[len(jokeDate)-2:])
		if joke.Rating > maxRate && today.Year() == year && int(today.Month()) == month && (today.Day()-day) == 1 {
			dailyJoke = joke
			maxRate = joke.Rating
		}
	}
	return dailyJoke
}

func (js *JokesStore) GetGeneratedJoke() Joke {
	return GetJokeFromDB(js.db, js.GeneratedJokeId)
}
