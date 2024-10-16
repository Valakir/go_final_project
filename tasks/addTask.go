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
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// Парсим JSON
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Ошибка чтения JSON", http.StatusBadRequest)
		return
	}
	// Проверка на пустые поля
	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка даты
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dates.DateFormat)
	} else {
		_, err := time.Parse(dates.DateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неправильный формат даты"}`, http.StatusBadRequest)
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
				http.Error(w, `{"error":"Некорректное значение даты повторения"}`, http.StatusBadRequest)
				return
			}
		}
	}

	// Проверка на повторение
	taskDate, err := time.Parse(dates.DateFormat, task.Date)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"ошибка":"неправильный формат даты: %s"}`, err.Error()), http.StatusBadRequest)
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
				http.Error(w, fmt.Sprintf(`{"ошибка получения даты":"%s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	// Вставка записи в базу данных
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Ошибка вставки в базу данных"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"Ошибка получения ID"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"id": int(id)})

}
