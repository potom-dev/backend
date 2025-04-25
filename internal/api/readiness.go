package api

import "net/http"

// HandlerReadiness godoc
//
//	@Router		/readiness [get]
//	@Summary	check if the service is ready
//	@Tags		health
//	@Accept		json
//	@Produce	text/plain
//	@Success	200	{string}	string
func (cfg *Config) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
