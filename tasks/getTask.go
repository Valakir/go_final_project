package tasks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// GetTasksHandler Получение списка задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Устанавливаем тип контента
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Проверка метода
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	var tasks []Task
	var err error

	if id != "" {
		// Поиск задачи по идентификатору
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
		var task Task
		err = db.QueryRow(query, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error":"Ошибка при запросе к базе данных"}`, http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	} else {
		// Запрос всех задач, лимит до 50
		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50"
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при запросе к базе данных"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Обработка результатов запроса
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				http.Error(w, `{"error":"Ошибка при чтении результатов"}`, http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		// Проверка на ошибки сканирования
		if err = rows.Err(); err != nil {
			http.Error(w, `{"error":"Ошибка при обработке результатов"}`, http.StatusInternalServerError)
			return
		}
	}
	// если нет задач, возвращаем пустой список
	if len(tasks) == 0 {
		tasks = []Task{}
	}

	// Отправка JSON-ответа
	jsonResponse, err := json.Marshal(map[string]interface{}{"tasks": tasks})
	if err != nil {
		http.Error(w, `{"error":"Ошибка при кодировании ответа"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

// Поучение задачи для редактирования
func GetTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Получение идентификатора задачи из запроса
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// SQL-запрос для получения задачи по ID
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, taskID)

	var task Task
	var id int
	err := row.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка поиска задачи: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Преобразование идентификатора в строку
	task.Id = strconv.Itoa(id)

	// Отправка задачи в формате JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
