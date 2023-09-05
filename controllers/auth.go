package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/TianMeh/go-guest/models"
	"github.com/TianMeh/go-guest/utils"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required"`
}

func Signup(w http.ResponseWriter, r *http.Request) {

	var creds Credentials

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &creds)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		log.Fatal("Error unmarshalling signup request data:", err.Error())
		return
	}

	validate = validator.New()
	err = validate.Struct(creds)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, string(err.Error()))
		log.Fatal("Invalid data:", err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, string(err.Error()))
		log.Fatal("Error hashing password:", err.Error())
		return
	}

	user := &models.User{
		Password: string(hashedPassword),
		Username: creds.Username,
	}

	result := models.DB.Where("username = ?", user.Username).First(&models.User{})

	if result.RowsAffected > 0 {
		// User with the desired username already exists
		utils.RespondWithError(w, http.StatusUnauthorized, "User already exists.")
		return
	}

	models.DB.Create(user)

	setHeader(w)
	json.NewEncoder(w).Encode(user)

}

func Signin(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, creds)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, string(err.Error()))
		log.Fatal("Error unmarshalling request body", err.Error())
	}

	storedCreds := &Credentials{}
	result := models.DB.Where("username = ?", creds.Username).First(&models.User{}).Scan(storedCreds)

	if result.RowsAffected == 0 {
		// User with the desired username already exists
		utils.RespondWithError(w, http.StatusNotFound, "No matching user found")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Wrong credentials")
		return
	}

	setHeader(w)
	response := map[string]string{"username": storedCreds.Username}
	json.NewEncoder(w).Encode(response)

	return

}
