package tasks

import (
	"database/sql"
	"fmt"
	"net/http"
)

// DeleteTaskHandler обрабатывает DELETE запрос для удаления задачи по ID.
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Получаем идентификатор задачи из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		RespondWithError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	// SQL-запрос на удаление задачи
	result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка удаления задачи: %s", err.Error()))
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка получения результата операции: %s", err.Error()))
		return
	}

	// Проверка, была ли удалена запись
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	// Возврат успешного ответа
	RespondWithSuccess(w)
}
