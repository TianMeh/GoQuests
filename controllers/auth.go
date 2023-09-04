package controllers

import (
	"encoding/json"
	"io"
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

	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &creds)

	validate = validator.New()
	err := validate.Struct(creds)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, string(err.Error()))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	user := &models.User{
		Password: string(hashedPassword),
		Username: creds.Username,
	}

	models.DB.Create(user)

	setHeader(w)
	json.NewEncoder(w).Encode(user)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	return
}
