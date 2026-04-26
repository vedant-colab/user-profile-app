package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, ErrorResponse{
		Error: msg,
	})
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{
		Status: "success",
		Data:   data,
	})
}
