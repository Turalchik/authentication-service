package handlers

import "net/http"

func (httpHandler *HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	httpHandler.router.ServeHTTP(w, req)
}
