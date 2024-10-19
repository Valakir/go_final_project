package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/dates"
	"net/http"
	"strings"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Устанавливаем тип контента
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Проверка метода
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}
	// Парсим JSON
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Ошибка чтения JSON")
		return
	}
	// Проверка на пустые поля
	if task.Title == "" {
		RespondWithError(w, http.StatusBadRequest, `{"error":"Не указан заголовок задачи"}`)
		return
	}

	// Проверка даты
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dates.DateFormat)
	} else {
		_, err := time.Parse(dates.DateFormat, task.Date)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, `{"error":"Неправильный формат даты"}`)
			return
		}
	}
	// Валидация поля Repeat
	validRepeats := map[string]bool{
		"d": true, "y": true,
	}
	if task.Repeat != "" {
		parts := strings.Fields(task.Repeat)
		if len(parts) > 0 {
			_, isValid := validRepeats[parts[0]]
			if !isValid {
				RespondWithError(w, http.StatusBadRequest, `{"error":"Некорректное значение даты повторения"}`)
				return
			}
		}
	}

	// Проверка на повторение
	taskDate, err := time.Parse(dates.DateFormat, task.Date)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"ошибка":"неправильный формат даты: %s"}`, err.Error()))
		return
	}

	// Обрезаем время до даты без времени
	taskDate, nowDate := taskDate.Truncate(24*time.Hour), now.Truncate(24*time.Hour)

	if taskDate.Before(nowDate) {
		if task.Repeat == "" {
			// Если нет правила повторения, обновляем дату до текущей
			task.Date = nowDate.Format(dates.DateFormat)
		} else {
			// Если есть правило повторения, определяем следующую дату
			nextDate, err := dates.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"ошибка получения даты":"%s"}`, err.Error()))
				return
			}
			task.Date = nextDate
		}
	}

	// Вставка записи в базу данных
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, `{"error":"Ошибка вставки в базу данных"}`)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, `{"error":"Ошибка получения ID"}`)
		return
	}

	// Отправка задачи в формате JSON
	json.NewEncoder(w).Encode(map[string]int{"id": int(id)})

}
