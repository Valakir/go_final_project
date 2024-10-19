package tasks

import (
	"encoding/json"
	"net/http"
)

// Task - структура задачи
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// JSONResponse содержит формат ответа
type JSONResponse map[string]interface{}

// RespondWithError отправляет JSON с ошибкой
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RespondWithSuccess отправляет успешный пустой JSON
func RespondWithSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSONResponse{})
}
