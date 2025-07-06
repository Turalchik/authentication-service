package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (authHandler *HttpHandler) UpdateTokens(w http.ResponseWriter, req *http.Request) {
	body := tokensBody{}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	userAgent := req.UserAgent()
	ipAddr, _ := getIP(req)

	accessToken, refreshToken, err := authHandler.authService.UpdateTokens(body.AccessToken, body.RefreshToken, userAgent, ipAddr)
	if err != nil {
		// TODO
		// распарсить ошибки
		http.Error(w, fmt.Sprintf("Invalid request: %s", err.Error()), http.StatusBadRequest)
	}

	resp := tokensBody{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	var jsonBytes []byte
	if jsonBytes, err = json.MarshalIndent(resp, "", "\t"); err != nil {
		http.Error(w, fmt.Sprintf("Can't marshal response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(jsonBytes); err != nil {
		http.Error(w, fmt.Sprintf("Can't send response: %s", err.Error()), http.StatusInternalServerError)
	}
}
