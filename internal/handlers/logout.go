package handlers

import (
	"net/http"
)

// Logout разлогинивает текущего пользователя, отзывая его refresh‑токен.
// @Summary      Выход пользователя (logout)
// @Description  Инвалидирует refresh‑токен текущего пользователя, после чего refresh и protected‑маршруты недоступны.
// @Tags         auth
// @Security     ApiKeyAuth
// @Success      204  {string}  string  "No Content"
// @Failure      400  {string}  string  "invalid access token"
// @Failure      401  {string}  string  "unauthorized"
// @Router       /api/v1/auth/logout [post]
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
