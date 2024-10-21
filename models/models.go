package models

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


// GetTaskIDFromRequest получает идентификатор задачи из параметров запроса.
func GetTaskIDFromRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := r.URL.Query().Get("id")
	if id == "" {
		RespondWithError(w, http.StatusBadRequest, "Не указан идентификатор")
		return "", false
	}
	return id, true
}

// ValidateHTTPMethod проверяет, соответствует ли метод запроса ожидаемому.
func ValidateHTTPMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		RespondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return false
	}
	return true
}
