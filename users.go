package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users", err)
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	user, err := cfg.db.GetUser(r.Context(), uuid.MustParse(userId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithError(w, http.StatusMethodNotAllowed, "Not allowed", nil)
		return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}
	respondWithJSON(w, http.StatusOK, "All users deleted")
}
