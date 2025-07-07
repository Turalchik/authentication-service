package handlers

import (
	"encoding/json"
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

	resp := tokensBody{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	var jsonBytes []byte
	if jsonBytes, err = json.MarshalIndent(resp, "", "\t"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(jsonBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
