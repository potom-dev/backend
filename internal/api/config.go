package api

import (
	"net/http"
	"sync/atomic"

	"github.com/potom-dev/backend/internal/database"
)

type Config struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	jwtSecret      string
}

func NewConfig(db *database.Queries, jwtSecret string) *Config {
	return &Config{
		fileserverHits: atomic.Int32{},
		db:             db,
		jwtSecret:      jwtSecret,
	}
}

func (cfg *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
