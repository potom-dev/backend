package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           potom swagger api
// @version         1.0
// @description     potom api.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func NewRouter(cfg *Config) http.Handler {
	mux := http.NewServeMux()

	fsHandler := cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", cfg.HandlerReadiness)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /api/users", cfg.handlerGetUsers)
	mux.HandleFunc("GET /api/users/{userId}", cfg.handlerGetUser)
	mux.HandleFunc("PUT /api/users/{userId}", cfg.handlerUpdateUser)
	mux.HandleFunc("DELETE /api/users", cfg.handlerDeleteAllUsers)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeRefresh)

	mux.HandleFunc("POST /api/groups", cfg.handlerCreateGroup)

	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	return mux
}
