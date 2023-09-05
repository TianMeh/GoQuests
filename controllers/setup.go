package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func New() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/quests", GetAllQuests).Methods("GET")
	router.HandleFunc("/quest/{id}", GetQuest).Methods("GET")
	router.HandleFunc("/quest", CreateQuest).Methods("POST")
	router.HandleFunc("/quest/{id}", UpdateQuest).Methods("PUT")
	router.HandleFunc("/quest/{id}", DeleteQuest).Methods("DELETE")

	// Auth
	router.HandleFunc("/signup", Signup).Methods("POST")
	router.HandleFunc("/signin", Signin).Methods("POST")
	router.HandleFunc("/signout", Signout).Methods("POST")

	return router
}
