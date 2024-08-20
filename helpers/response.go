package helpers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ResponseWithData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Response(w http.ResponseWriter, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	status := code <= 400

	response := ResponseWithData{
		Status:  strconv.FormatBool(status),
		Message: message,
		Data:    data,
	}

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		panic(err)
	}
}
