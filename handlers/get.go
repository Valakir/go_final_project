package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go_final_project/models"
)

const tasksLimit = 50

// GetTasksHandler Получение списка задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Устанавливаем тип контента
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Проверка метода
	if !models.ValidateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	id := r.URL.Query().Get("id")
	var tasks []models.Task
	var err error

	if id != "" {
		// Поиск задачи по идентификатору
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
		var task models.Task
		err = db.QueryRow(query, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if errors.Is(err, sql.ErrNoRows) {
			models.RespondWithError(w, http.StatusNotFound, `{"error":"Задача не найдена"}`)
			return
		} else if err != nil {
			models.RespondWithError(w, http.StatusInternalServerError, `{"error":"Ошибка при запросе к базе данных"}`)
			return
		}
		tasks = append(tasks, task)
	} else {
		// Запрос всех задач, лимит до 50
		query := fmt.Sprintf("SELECT id, date, title, comment, repeat FROM scheduler LIMIT %d", tasksLimit)
		rows, err := db.Query(query)
		if err != nil {
			models.RespondWithError(w, http.StatusInternalServerError, `{"error":"Ошибка при запросе к базе данных"}`)
			return
		}
		defer rows.Close()

		// Обработка результатов запроса
		for rows.Next() {
			var task models.Task
			if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				models.RespondWithError(w, http.StatusInternalServerError, `{"error":"Ошибка при чтении результатов"}`)
				return
			}
			tasks = append(tasks, task)
		}

		// Проверка на ошибки сканирования
		if err = rows.Err(); err != nil {
			models.RespondWithError(w, http.StatusInternalServerError, `{"error":"Ошибка при обработке результатов"}`)
			return
		}
	}
	// если нет задач, возвращаем пустой список
	if len(tasks) == 0 {
		tasks = []models.Task{}
	}

	// Отправка JSON-ответа
	jsonResponse, err := json.Marshal(map[string]interface{}{"tasks": tasks})
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, `{"error":"Ошибка при кодировании ответа"}`)
		return
	}

	w.Write(jsonResponse)
}

// Поучение задачи для редактирования
func GetTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Получение идентификатора из запроса
	taskID, ok := models.GetTaskIDFromRequest(w, r)
	if !ok {
		return
	}

	// SQL-запрос для получения задачи по ID
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, taskID)

	var task models.Task
	var id int
	err := row.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		models.RespondWithError(w, http.StatusNotFound, `{"error":"Задача не найдена"}`)
		return
	} else if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf(`{"error":"Ошибка поиска задачи: %s"}`, err.Error()))
		return
	}

	// Преобразование идентификатора в строку
	task.Id = strconv.Itoa(id)

	// Отправка задачи в формате JSON

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
