package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
		http.Error(w, fmt.Sprintf("Invalid request: %s", err.Error()), http.StatusBadRequest)
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
