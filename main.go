package main

import (
	"net/http"

	"github.com/TianMeh/go-guest/controllers"
	"github.com/TianMeh/go-guest/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	handler := controllers.New()

	server := &http.Server{
		Addr:    "0.0.0.0:8008",
		Handler: handler,
	}

	models.ConnectDatabase()

	server.ListenAndServe()

}
