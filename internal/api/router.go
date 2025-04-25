package api

import (
	"net/http"
)

func NewRouter(cfg *Config) http.Handler {
	mux := http.NewServeMux()

	fsHandler := cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", cfg.HandlerReadiness)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /api/users", cfg.handlerGetUsers)
	mux.HandleFunc("GET /api/users/{userId}", cfg.handlerGetUser)
	mux.HandleFunc("DELETE /api/users", cfg.handlerDeleteAllUsers)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeRefresh)

	mux.HandleFunc("POST /api/groups", cfg.handlerCreateGroup)

	return mux
}
