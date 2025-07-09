package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Guid возвращает GUID текущего пользователя.
// @Summary      Получение GUID текущего пользователя
// @Description  Возвращает GUID пользователя, извлечённый из access token.
// @Tags         auth
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  userIDBody
// @Failure      401  {string}  string  "unauthorized"
// @Router       /api/v1/auth/guid [get]
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
