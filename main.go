package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/potom-dev/backend/internal/api"
	"github.com/potom-dev/backend/internal/auth"
	"github.com/potom-dev/backend/internal/database"
	"github.com/potom-dev/backend/internal/env"

	// Import pq driver for its side effects only
	_ "github.com/lib/pq"
	_ "github.com/potom-dev/backend/docs"
)

func main() {
	env.InitEnv()

	auth.NewAuth()

	port := env.GetEnv("PORT")
	jwtSecret := env.GetEnv("JWT_SECRET")
	dbURL := env.GetEnv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	apiCfg := api.NewConfig(dbQueries, jwtSecret)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: api.NewRouter(apiCfg),
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
