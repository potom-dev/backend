package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/potom-dev/backend/internal/auth"
	"github.com/potom-dev/backend/internal/database"
)

type CreateUpdateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// handlerCreateUser godoc
//
//	@Router		/users [post]
//	@Summary	create a user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		body	body		CreateUpdateUserParams	true	"User creation parameters"
//	@Success	201		{object}	User
//	@Failure	500		{object}	ErrorResponse
func (cfg *Config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := CreateUpdateUserParams{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	pswdHash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	log.Printf("Hashed password: %s", pswdHash)

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:        params.Email,
		PasswordHash: pswdHash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// handlerGetUsers godoc
//
//	@Router		/users [get]
//	@Summary	get all users
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		User
//	@Failure	500	{object}	ErrorResponse
func (cfg *Config) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users", err)
		return
	}

	usersResponse := []User{}

	for _, user := range users {
		usersResponse = append(usersResponse, User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, usersResponse)
}

// handlerGetUser godoc
//
//	@Router		/users/{userId} [get]
//	@Summary	get user by id
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		userId	path		string	true	"User ID"
//	@Success	200		{object}	User
//	@Failure	500		{object}	ErrorResponse
func (cfg *Config) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	user, err := cfg.db.GetUserById(r.Context(), uuid.MustParse(userId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// handlerUpdateUser godoc
//
//	@Router		/users/{userId} [put]
//	@Summary	update user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		userId	path	string				true	"User ID"
//	@Param		body	body	CreateUpdateUserParams	true	"User update data"
//	@Success	204		"No Content"
//	@Failure	401		{object}	ErrorResponse
//	@Failure	403		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Security	BearerAuth
func (cfg *Config) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	authedUserID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	if authedUserID != uuid.MustParse(userId) {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	user, err := cfg.db.GetUserById(r.Context(), uuid.MustParse(userId))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := CreateUpdateUserParams{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	pswdHash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:           user.ID,
		Email:        params.Email,
		PasswordHash: pswdHash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}

// handlerDeleteAllUsers godoc
//
//	@Router		/users [delete]
//	@Summary	delete all users
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{string}	string
//	@Failure	405	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
func (cfg *Config) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
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
