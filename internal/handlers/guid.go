package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (httpHandler *HttpHandler) Guid(w http.ResponseWriter, req *http.Request) {
	args := req.Context().Value("args").(map[string]string)

	resp := &userIDBody{
		UserID: args["userID"],
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Guid: failed to write response: %v", err)
	}
}
