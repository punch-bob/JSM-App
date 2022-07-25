package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	DB_NAME     = os.Getenv("DB_NAME")
	DB_LOGIN    = os.Getenv("DB_LOGIN")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
)

type Joke struct {
	Id         int      `json:"id"`
	Text       string   `json:"text"`
	Rating     int      `json:"rate"`
	Tags       []string `json:"tags"`
	AuthorName string   `json:"author_name"`
	Date       string   `json:"date"`
	UserId     int      `json:"uid"`
}

type JokesStore struct {
	GeneratedJoke Joke
	db            *sql.DB
}

func newJokesStore() *JokesStore {
	GeneratedJoke := Joke{}
	dataSourceName := DB_LOGIN + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	return &JokesStore{
		GeneratedJoke: GeneratedJoke,
		db:            db,
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

func (js *JokesStore) getJokeTags(id int) []string {
	var tag []string
	tags, err := js.db.Query("select tag from "+DB_NAME+".tag where joke_id = ?", id)
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

func (js *JokesStore) getJokeById(id int) Joke {
	row := js.db.QueryRow("select * from "+DB_NAME+".joke where id = ?", id)
	joke := Joke{}
	err := row.Scan(&joke.Id, &joke.UserId, &joke.Text, &joke.Rating, &joke.AuthorName, &joke.Date)
	if err != nil {
		log.Println(err)
	}
	joke.Tags = js.getJokeTags(id)
	return joke
}

func (js *JokesStore) convertRowToJoke(row *sql.Rows) Joke {
	joke := Joke{}
	err := row.Scan(&joke.Id, &joke.UserId, &joke.Text, &joke.Rating, &joke.AuthorName, &joke.Date)
	if err != nil {
		log.Println(err)
	}
	joke.Tags = js.getJokeTags(joke.Id)
	return joke
}

func (js *JokesStore) getJokeList() []Joke {
	var jokeStore []Joke
	rows, err := js.db.Query("select * from " + DB_NAME + ".joke")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		joke := Joke{}
		err := rows.Scan(&joke.Id, &joke.UserId, &joke.Text, &joke.Rating, &joke.AuthorName, &joke.Date)
		if err != nil {
			continue
		}
		joke.Tags = js.getJokeTags(joke.Id)
		jokeStore = append(jokeStore, joke)
	}
	return jokeStore
}

func (js *JokesStore) updateRating(id, rating int) {
	updateRating, err := js.db.Prepare("update joke set rating = ? where id = ?")
	if err != nil {
		log.Println(err)
	}
	_, _ = updateRating.Exec(rating, id)
}

func (js *JokesStore) updateReact(uid, id, react int) {
	updateReact, err := js.db.Prepare("update joke_rating set react = ? where joke_id = ? and user_id = ?")
	if err != nil {
		log.Println(err)
	}
	_, _ = updateReact.Exec(react, id, uid)
}

func (js *JokesStore) deleteJoke(id int) {
	_, err := js.db.Exec("delete from "+DB_NAME+".joke where id = ?", id)
	if err != nil {
		log.Println(err)
	}

	_, err = js.db.Exec("delete from "+DB_NAME+".joke_rating where joke_id = ?", id)
	if err != nil {
		log.Println(err)
	}

	_, err = js.db.Exec("delete from "+DB_NAME+".tag where joke_id = ?", id)
	if err != nil {
		log.Println(err)
	}
}

func (js *JokesStore) addJoke(joke Joke) int {
	var (
		text       = joke.Text
		rating     = joke.Rating
		authorName = joke.AuthorName
		date       = joke.Date
		uid        = joke.UserId
	)
	res, err := js.db.Exec("insert into "+DB_NAME+".joke (uid, text, rating, author_name, date) values (?, ?, ?, ?, ?)", uid, text, rating, authorName, date)
	if err != nil {
		log.Println(err)
		return -1
	}
	id, _ := res.LastInsertId()
	joke.Id = int(id)

	for _, tag := range joke.Tags {
		_, err := js.db.Exec("insert into "+DB_NAME+".tag (joke_id, tag) values (?, ?)", id, tag)
		if err != nil {
			log.Println(err)
		}
	}
	return int(id)
}

func (js *JokesStore) createJoke(text string, tags []string, authorName string, uid int) Joke {
	joke := Joke{
		Id:         0,
		Text:       text,
		Rating:     0,
		AuthorName: authorName,
		Date:       time.Now().Format("2006-01-02"),
		UserId:     uid}
	joke.Tags = make([]string, len(tags))
	copy(joke.Tags, tags)

	js.addJoke(joke)

	return joke
}

func (js *JokesStore) increaseRating(uid, id int) int {
	joke := js.getJokeById(id)

	stmt, err := js.db.Prepare("select react from " + DB_NAME + ".joke_rating where joke_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
	}
	react := stmt.QueryRow(id, uid)
	value := 1
	err = react.Scan(&value)
	if err == sql.ErrNoRows {
		js.db.Exec("insert into "+DB_NAME+".joke_rating (joke_id, user_id, react) values (?, ?, ?)", id, uid, value)
	}

	if value == 2 {
		return joke.Rating
	}

	joke.Rating++
	js.updateRating(id, joke.Rating)
	js.updateReact(uid, id, value+1)
	return joke.Rating
}

func (js *JokesStore) decreaseRating(uid, id int) int {
	joke := js.getJokeById(id)

	stmt, err := js.db.Prepare("select react from " + DB_NAME + ".joke_rating where joke_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
	}
	react := stmt.QueryRow(id, uid)
	value := 1
	err = react.Scan(&value)
	if err == sql.ErrNoRows {
		js.db.Exec("insert into "+DB_NAME+".joke_rating (joke_id, user_id, react) values (?, ?, ?)", id, uid, value)
	}

	if value == 0 {
		return joke.Rating
	}

	joke.Rating--
	js.updateRating(id, joke.Rating)
	js.updateReact(uid, id, value-1)
	return joke.Rating
}

func (js *JokesStore) getJokeIdByTag(tag string) []int {
	var jokeId []int

	jokes, err := js.db.Query("select joke_id from "+DB_NAME+".tag where tag = ?", tag)
	if err != nil {
		log.Println(err)
	}
	defer jokes.Close()

	for jokes.Next() {
		var id int
		err = jokes.Scan(&id)
		if err != nil {
			continue
		}

		jokeId = append(jokeId, id)
	}
	return jokeId
}

func (js *JokesStore) getJokesByTagList(tags []string) []Joke {
	var jokes []Joke
	var jokeId []int

	for _, tag := range tags {
		jokeId = append(jokeId, js.getJokeIdByTag(tag)...)
		jokeId = removeDuplicateValues(jokeId)
	}

	for _, id := range jokeId {
		jokes = append(jokes, js.getJokeById(id))
	}
	return jokes
}

func (js *JokesStore) getDailyJoke() Joke {
	var dailyJoke Joke

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	const MaxInt = int(^uint(0) >> 1)
	maxRate := -MaxInt - 1

	jokes, err := js.db.Query("select * from " + DB_NAME + ".joke")
	if err != nil {
		if err == sql.ErrNoRows {
			return Joke{}
		}
		log.Println(err)
	}
	defer jokes.Close()

	for jokes.Next() {
		joke := js.convertRowToJoke(jokes)
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

func (js *JokesStore) getJokesByUID(uid int) []Joke {
	var jokeList []Joke

	jokes, err := js.db.Query("select * from "+DB_NAME+".joke where uid = ?", uid)
	if err != nil {
		log.Println(err)
	}
	defer jokes.Close()

	for jokes.Next() {
		joke := js.convertRowToJoke(jokes)
		jokeList = append(jokeList, joke)
	}
	return jokeList
}
