package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/dates"
	"net/http"
	"time"
)

// UpdateTaskHandle обновляем запись задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Некорректный запрос"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация даты
	if _, err := time.Parse(dates.DateFormat, task.Date); err != nil {
		http.Error(w, `{"error":"Некорректная дата"}`, http.StatusBadRequest)
		return
	}

	// Проверяем что дата не в прошлом
	if task.Date < time.Now().Format(dates.DateFormat) {
		http.Error(w, `{"error":"Дата не может быть меньше сегодняшней"}`, http.StatusBadRequest)
		return
	}

	// Проверяем что заголовок заполнен
	if task.Title == "" {
		http.Error(w, `{"error":"Заголовок не может быть пустым"}`, http.StatusBadRequest)
		return
	}

	// SQL запрос для обновления
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка обновления задачи: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка получения результата обновления: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON если успех
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
