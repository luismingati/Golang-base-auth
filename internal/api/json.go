package api

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json %v\n", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, errMessage string, message string) {
	if code > 499 {
		slog.Error(fmt.Sprint("Responding with 5XX error: ", message))
	}
	type errResponse struct {
		Code    int    `json:"code"`
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	respondWithJSON(w, code, errResponse{Code: code, Error: errMessage, Message: message})
}

func respondWithInternalServerError(w http.ResponseWriter) {
	respondWithError(w, http.StatusInternalServerError, "Erro interno no servidor.", "Tente novamente mais tarde.")
}
