package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var regexpValidToken = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]*$")

type AuthServer struct {
	db *sql.DB
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	Id            int    `json:"id"`
	ServerMessage string `json:"server_message"`
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func checkUsernameValidation(username string) bool {
	return regexpValidToken.MatchString(username) && len(username) > 0
}

func checkPasswordValidation(password string) bool {
	return regexpValidToken.MatchString(password) && len(password) >= 7
}

func checkPasswordStrength(password string) int {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return boolToInt(hasUpper) + boolToInt(hasLower) + boolToInt(hasNumber) + boolToInt(hasSpecial)
}

func (server *AuthServer) CheckUserPassword(username, password string) (int, bool) {
	type TableRaw struct {
		Id       int
		Username string
		Password []byte
	}
	// rows, err := server.db.Query("select * from "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" where username = ?", username)
	rows, err := server.db.Query("select * from account_db.user where username = ?", username)
	if err != nil {
		log.Println(err)
		return -1, false
	}

	defer rows.Close()
	users := []TableRaw{}

	for rows.Next() {
		user := TableRaw{}
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}

	for _, user := range users {
		err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
		if err == nil {
			return user.Id, true
		}
	}
	return -1, false
}

func (server *AuthServer) AddUser(name, password string) (int, error) {
	if !checkPasswordValidation(password) {
		return -1, fmt.Errorf("your password is invalid")
	}

	passLevel := checkPasswordStrength(password)
	if passLevel <= 2 {
		return -1, fmt.Errorf("your password is weak: %d level of 4", passLevel)
	}

	if !checkUsernameValidation(name) {
		return -1, fmt.Errorf("your username is invalid")
	}

	tmp := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(tmp, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	// res, err := server.db.Exec("insert into "+os.Getenv("DB_NAME")+"."+os.Getenv("TABLE_NAME")+" (username, password) values (?, ?)", name, hashedPassword)
	res, err := server.db.Exec("insert into account_db.user (username, password) values (?, ?)", name, hashedPassword)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return int(id), nil
}

func (server *AuthServer) AuthorizationHandler(write http.ResponseWriter, request *http.Request) {
	dec := json.NewDecoder(request.Body)
	var user User
	err := dec.Decode(&user)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	var response Response
	id, err := server.AddUser(user.Name, user.Password)
	response.Id = id
	if err != nil {
		response.ServerMessage = err.Error()
	} else {
		response.ServerMessage = "Ok"
	}

	js, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *AuthServer) LogUpHandler(write http.ResponseWriter, request *http.Request) {
	dec := json.NewDecoder(request.Body)
	var user User
	err := dec.Decode(&user)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	var response Response
	id, correct := server.CheckUserPassword(user.Name, user.Password)
	response.Id = id
	if !correct {
		response.ServerMessage = "Wrong password or username"
	} else {
		response.ServerMessage = "Ok"
	}

	js, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func main() {
	log.Println("start")
	var server AuthServer
	var err error
	// dataSourceName := os.Getenv("DB_LOGIN") + ":" + os.Getenv("DB_PASSWORD") + "@/" + os.Getenv("DB_NAME")
	dataSourceName := "root:password@/account_db"
	server.db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	defer server.db.Close()
	router := mux.NewRouter()

	router.HandleFunc("/authorization/", server.AuthorizationHandler).Methods("POST")
	router.HandleFunc("/log_up/", server.LogUpHandler).Methods("POST")

	//Header sets
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "HEAD", "PUT", "OPTIONS"})

	// http.ListenAndServe(os.Getenv("SERVER_PORT"), handlers.CORS(headersOk, originsOk, methodsOk)(router))
	http.ListenAndServe(":8082", handlers.CORS(headersOk, originsOk, methodsOk)(router))
}
