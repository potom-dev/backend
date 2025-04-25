package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/potom-dev/backend/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		ExpiresInSec int64  `json:"expires_in_seconds,omitempty"`
	}
	type response struct {
		Id    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		Token string    `json:"token"`
	}

	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPassword(params.Password, user.PasswordHash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(params.ExpiresInSec)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Id:    user.ID,
		Email: user.Email,
		Token: token,
	})
}
