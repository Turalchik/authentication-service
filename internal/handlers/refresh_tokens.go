package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RefreshTokens обновляет пару токенов по существующему access + refresh.
// @Summary      Обновление токенов
// @Description  Принимает действующий access и refresh‑токены, возвращает новую пару.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      accessAndRefreshTokensBody  true  "Существующие access и refresh"
// @Success      200   {object}  accessAndRefreshTokensBody
// @Failure      400   {string}  string  "Invalid request body or tokens"
// @Failure      401   {string}  string  "Unauthorized or token revoked"
// @Router       /api/v1/auth/tokens/refresh [post]
func (httpHandler *HttpHandler) RefreshTokens(w http.ResponseWriter, req *http.Request) {
	body := accessAndRefreshTokensBody{}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	userAgent := req.UserAgent()
	ipAddr, _ := getIP(req)

	accessToken, refreshToken, err := httpHandler.authService.RefreshTokens(body.AccessToken, body.RefreshToken, userAgent, ipAddr)
	if err != nil {
		// TODO
		// распарсить ошибки
		http.Error(w, fmt.Sprintf("Invalid request: %s", err.Error()), http.StatusUnauthorized)
	}

	resp := &accessAndRefreshTokensBody{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("CreateTokens: failed to write response: %v", err)
	}
}
