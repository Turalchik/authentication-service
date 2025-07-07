package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

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
		w.WriteHeader(http.StatusInternalServerError)
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
