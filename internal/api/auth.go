package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/potom-dev/backend/internal/auth"
	"github.com/potom-dev/backend/internal/database"
)

type LoginParams struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ExpiresInSec int64  `json:"expires_in_seconds,omitempty"`
}

type LoginResponse struct {
	Id           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type RefreshResponse struct {
	Token string `json:"token"`
}

// handlerLogin godoc
//
//	@Router		/login [post]
//	@Summary	login user
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		LoginParams	true	"Login parameters"
//	@Success	200		{object}	LoginResponse
//	@Failure	401		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
func (cfg *Config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var params LoginParams
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

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	refresh, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		Id:           user.ID,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

// handlerRefresh godoc
//
//	@Router		/refresh [post]
//	@Summary	refresh access token
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Success	200		{object}	RefreshResponse
//	@Failure	401		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
func (cfg *Config) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), refresh)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	user, err := cfg.db.GetUserById(r.Context(), refreshToken.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshResponse{
		Token: token,
	})
}

// handlerRevokeRefresh godoc
//
//	@Router		/refresh [delete]
//	@Summary	revoke refresh token
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Success	204		"No Content"
//	@Failure	401		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
func (cfg *Config) handlerRevokeRefresh(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refresh)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
