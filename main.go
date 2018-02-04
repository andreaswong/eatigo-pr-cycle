package eatigo_pr_cycle

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	http.ListenAndServe(":" + port, &DefaultHandler{})
}

type DefaultHandler struct {}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Query().Get("s")))
}

