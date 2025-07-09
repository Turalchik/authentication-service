package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CreateTokens выдаёт новую пару токенов (access + refresh).
// @Summary      Выдача токенов
// @Description  Генерирует пару токенов для пользователя с указанным user_id в query‑параметре.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user_id  query     string  true  "GUID пользователя"
// @Success      200      {object}  accessAndRefreshTokensBody
// @Failure      400      {string}  string  "user_id required"
// @Failure      500      {string}  string  "internal server error"
// @Router       /api/v1/auth/tokens [get]
func (httpHandler *HttpHandler) CreateTokens(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	userAgent := req.UserAgent()
	ipAddr, _ := getIP(req)

	accessToken, refreshToken, err := httpHandler.authService.CreateTokens(userID, userAgent, ipAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't create tokens with error: %v", err), http.StatusInternalServerError)
		return
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
