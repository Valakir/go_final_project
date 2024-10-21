package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"go_final_project/dates"
	"net/http"
	"time"

	"go_final_project/models"
)

// DoneTaskResponse структура ответа
type DoneTaskResponse struct {
	Error string `json:"error,omitempty"`
}

// DoneTaskHandler обрабатывает HTTP-запросы для отметки задачи как завершенной
func DoneTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Получение идентификатора из запроса
	id, ok := models.GetTaskIDFromRequest(w, r)
	if !ok {
		return
	}

	var date, repeat string
	// SQL-запрос для получения задачи по ID
	err := db.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?", id).Scan(&date, &repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			models.RespondWithError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка поиска: %s", err.Error()))
		}
		return
	}

	// Проверка на повтор задачи
	if repeat == "" {
		// Удаляем не повторяющуюся задачу
		if _, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id); err != nil {
			models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка удаления: %s", err.Error()))
		} else {
			models.RespondWithSuccess(w)
		}
		return
	}

	// Вычисление следующей даты для повторяемой задачи
	nextDate, err := dates.NextDate(time.Now(), date, repeat)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка расчета даты: %s", err.Error()))
		return
	}

	// Обновление даты задачи
	if _, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id); err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка обновления даты: %s", err.Error()))
	} else {
		models.RespondWithSuccess(w)
	}
}
