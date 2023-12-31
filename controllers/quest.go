package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/TianMeh/go-guest/models"
	"github.com/TianMeh/go-guest/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func CheckSession(r *http.Request) (bool, error) {

	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil
		}
		return false, err
	}

	sessionToken := c.Value
	var session models.Session
	result := models.DB.Where("token = ?", sessionToken).First(&session)

	if result.Error != nil {
		return false, result.Error
	}

	userID := session.UserID
	var user models.User
	userResult := models.DB.First(&user, userID)

	if userResult.Error != nil {
		return false, userResult.Error
	}

	if session.Expires.Before(time.Now()) {
		models.DB.Delete(session)
		return false, nil
	}

	return true, nil
}

func GetAllQuests(w http.ResponseWriter, r *http.Request) {
	userAuthenticated, err := CheckSession(r)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if !userAuthenticated {
		utils.RespondWithError(w, http.StatusUnauthorized, "Token invalid, missing or expired")
		return
	}
	setHeader(w)

	var quests []models.Quest
	models.DB.Find(&quests)

	json.NewEncoder(w).Encode(quests)
}

func GetQuest(w http.ResponseWriter, r *http.Request) {
	userAuthenticated, err := CheckSession(r)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if !userAuthenticated {
		utils.RespondWithError(w, http.StatusUnauthorized, "Token invalid, missing or expired")
		return
	}
	setHeader(w)

	id := mux.Vars(r)["id"]
	var quest models.Quest

	if err := models.DB.Where("id = ?", id).First(&quest).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Quest not found")
		return
	}

	json.NewEncoder(w).Encode(quest)
}

var validate *validator.Validate

type QuestInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Reward      int    `json:"reward" validate:"required"`
}

func CreateQuest(w http.ResponseWriter, r *http.Request) {

	userAuthenticated, err := CheckSession(r)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if !userAuthenticated {
		utils.RespondWithError(w, http.StatusUnauthorized, "Token invalid, missing or expired")
		return
	}

	var input QuestInput

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &input)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		log.Fatal("Error unmarshalling signup request data:", err.Error())
	}

	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, string(err.Error()))
		return
	}

	quest := &models.Quest{
		Title:       input.Title,
		Description: input.Description,
		Reward:      input.Reward,
	}

	models.DB.Create(quest)

	setHeader(w)

	json.NewEncoder(w).Encode(quest)
}

func UpdateQuest(w http.ResponseWriter, r *http.Request) {

	userAuthenticated, err := CheckSession(r)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if !userAuthenticated {
		utils.RespondWithError(w, http.StatusUnauthorized, "Token invalid, missing or expired")
		return
	}
	setHeader(w)

	id := mux.Vars(r)["id"]
	var quest models.Quest

	if err := models.DB.Where("id = ?", id).First(&quest).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Quest not found")
		return
	}

	var input QuestInput

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &input)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		log.Fatal("Error unmarshalling signup request data:", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, string(err.Error()))
		log.Fatal("Invalid data:", err.Error())
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation Error")
	}

	quest.Title = input.Title
	quest.Description = input.Description
	quest.Reward = input.Reward

	models.DB.Save(&quest)

	json.NewEncoder(w).Encode(quest)
}

func DeleteQuest(w http.ResponseWriter, r *http.Request) {

	userAuthenticated, err := CheckSession(r)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if !userAuthenticated {
		utils.RespondWithError(w, http.StatusUnauthorized, "Token invalid, missing or expired")
		return
	}
	setHeader(w)

	id := mux.Vars(r)["id"]

	var quest models.Quest

	if err := models.DB.Where("id = ?", id).First(&quest).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Quest not found")
		return
	}

	models.DB.Delete(&quest)

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(quest)
}
