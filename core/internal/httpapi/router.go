package httpapi

import "net/http"

func NewRouter(manager ProfileService) http.Handler {
	handler := NewProfileHandler(manager)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /api/profiles", handler.Create)
	mux.HandleFunc("GET /api/profiles/{id}", handler.Get)
	mux.HandleFunc("PATCH /api/profiles/{id}", handler.Update)
	return mux
}
