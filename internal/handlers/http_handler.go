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
	router.HandleFunc("/api/v1/auth/refresh", httpHandler.RefreshTokens).Methods(http.MethodPost)

	protectedRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	protectedRouter.Use(httpHandler.AuthMiddleware)
	protectedRouter.HandleFunc("/logout", httpHandler.Logout).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/guid", httpHandler.Guid).Methods(http.MethodGet)

	return httpHandler
}
