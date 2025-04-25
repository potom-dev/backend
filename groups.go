package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateGroup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	group, err := cfg.db.CreateGroup(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create group", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Id:        group.ID,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	})
}
