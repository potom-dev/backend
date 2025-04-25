package api

import (
	"net/http"

	"github.com/potom-dev/backend/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// gin-swagger middleware
// swagger embed files

// @title           potom swagger api
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func NewRouter(cfg *Config) http.Handler {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "potom ???"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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
