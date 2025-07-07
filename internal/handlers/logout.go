package handlers

import (
	"net/http"
)

func (httpHandler *HttpHandler) Logout(w http.ResponseWriter, req *http.Request) {
	args := req.Context().Value("args").(map[string]string)

	if err := httpHandler.authService.Logout(args["accessToken"], args["userID"]); err != nil {
		// TODO
		// тут тоже нужно распарсить ошибки дружище
		http.Error(w, "invalid access token", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
