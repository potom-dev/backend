package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/potom-dev/backend/internal/auth"
	"github.com/potom-dev/backend/internal/database"
)

type CreateGroupParams struct {
	Name string `json:"name"`
}

type Group struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// handlerCreateGroup godoc
//
//	@Router		/groups [post]
//	@Summary	create a group
//	@Tags		groups
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Param		body	body		CreateGroupParams	true	"Group creation parameters"
//	@Success	201		{object}	Group
//	@Failure	401		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Security	BearerAuth
func (cfg *Config) handlerCreateGroup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := CreateGroupParams{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	group, err := cfg.db.CreateGroup(r.Context(), database.CreateGroupParams{
		Name:     params.Name,
		AuthorID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create group", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Group{
		Id:        group.ID,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	})
}
