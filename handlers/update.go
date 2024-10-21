package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/dates"
	"go_final_project/models"
)

// UpdateTaskHandle обновляем запись задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task models.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		models.RespondWithError(w,
			http.StatusBadRequest, `{"error":"Некорректный запрос"}`)
		return
	}
	defer r.Body.Close()

	// Валидация даты
	if _, err := time.Parse(dates.DateFormat, task.Date); err != nil {
		models.RespondWithError(w, http.StatusBadRequest, `{"error":"Некорректная дата"}`)
		return
	}

	// Проверяем что дата не в прошлом
	if task.Date < time.Now().Format(dates.DateFormat) {
		models.RespondWithError(w, http.StatusBadRequest, `{"error":"Дата не может быть меньше сегодняшней"}`)
		return
	}

	// Проверяем что заголовок заполнен
	if task.Title == "" {
		models.RespondWithError(w, http.StatusBadRequest, `{"error":"Заголовок не может быть пустым"}`)
		return
	}

	// SQL запрос для обновления
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf(`{"error":"Ошибка обновления задачи: %s"}`, err.Error()))
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf(`{"error":"Ошибка получения результата обновления: %s"}`, err.Error()))
		return
	}

	if rowsAffected == 0 {
		models.RespondWithError(w, http.StatusNotFound, `{"error":"Задача не найдена"}`)
		return
	}

	// Возвращаем пустой JSON если успех

	models.RespondWithSuccess(w)
}
