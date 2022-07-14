package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type AuthServer struct{}

func (server *AuthServer) AuthorizationHandler(write http.ResponseWriter, request *http.Request) {

	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *AuthServer) UpdateJokeRatingHandler(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func (server *AuthServer) Handler(write http.ResponseWriter, request *http.Request) {
	write.Header().Set("Access-Control-Allow-Origin", "*")
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}

func main() {
	router := mux.NewRouter()
	var server AuthServer

	router.HandleFunc("/authorization/", server.AuthorizationHandler).Methods("POST")
	router.HandleFunc("/update_rating/", server.UpdateJokeRatingHandler).Methods("POST")

	//Header sets
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "HEAD", "PUT", "OPTIONS"})

	http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(router))
}
