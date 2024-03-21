package http

import (
	"encoding/json"
	"net/http"
)

type errBody struct {
	Error string `json:"error"`
}

func J(w http.ResponseWriter, body any) {
	w.Header().Add("Content-Type", "application/json")
	err, ok := body.(error)
	if ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errBody{Error: err.Error()})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	}
}
