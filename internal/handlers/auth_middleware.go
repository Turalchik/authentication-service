package handlers

import (
	"context"
	"net/http"
)

func (httpHandler *HttpHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if len(auth) < 7 || auth[:7] != "Bearer " {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		// просим сервис проверить токен за нас
		tokenStr := auth[7:]
		userID, err := httpHandler.authService.CheckAccessTokenValidity(tokenStr)
		if err != nil {
			// TODO
			// как же много нужно парсить ошибки
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// прокидываем userID и access токен в контекст
		ctx := context.WithValue(req.Context(), "args", map[string]string{
			"userID":      userID,
			"accessToken": tokenStr,
		})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
