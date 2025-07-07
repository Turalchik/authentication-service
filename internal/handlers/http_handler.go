package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type HttpHandler struct {
	authService AuthService
	router      *mux.Router
}

func NewHttpHandler(authService AuthService) *HttpHandler {
	router := mux.NewRouter()
	httpHandler := &HttpHandler{
		authService: authService,
		router:      router,
	}

	router.HandleFunc("/api/v1/auth/tokens?user_id={GUID}", httpHandler.CreateTokens).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/auth/tokens/refresh", httpHandler.RefreshTokens).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/auth/tokens/logout", httpHandler.Logout).Methods(http.MethodPost)

	return httpHandler
}
