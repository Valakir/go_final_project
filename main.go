package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"go_final_project/auth"
	"go_final_project/database"
	"go_final_project/dates"
	"go_final_project/tasks"
)

// main это точка входа в приложение.
// Инициализирует базу данных, устанавливает маршруты и запускает HTTP-сервер.
func main() {
	// Инициализация БД
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("ошибка при подключении к БД: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("ошибка при закрытии БД: %v", err)
		}
	}()
	log.Println("БД подключена")

	// Установка маршрутов файлового сервера и API
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", dates.ApiNextDateHandler)

	// Создание обработчика задач
	taskHandler := &TaskHandler{db: db}
	http.HandleFunc("/api/tasks", taskHandler.GetTasksHandler)
	http.HandleFunc("/api/task/done", taskHandler.DoneTaskHandler)
	http.HandleFunc("/api/task", taskHandler.RouteTaskMethods)
	http.HandleFunc("/api/signin", auth.AuthUser(taskHandler.SignInHandler))

	// Получение порта из переменной окружения TODO_PORT или значение по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Запуск сервера и прослушивание порта
	log.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}

// TaskHandler является структурой, содержащей указатель на базу данных.
// Используется для передачи доступа к БД обработчикам задач.
type TaskHandler struct {
	db *sql.DB
}

// GetTasksHandler обрабатывает HTTP-запросы для получения списка задач.
func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks.GetTasksHandler(w, r, h.db)
}

// DoneTaskHandler обрабатывает HTTP-запросы для отметки задачи как завершенной.
func (h *TaskHandler) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks.DoneTaskHandler(w, r, h.db)
}

// SignInHandler обрабатывает HTTP-запросы для входа в систему.
func (h *TaskHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	auth.SignInHandler(w, r)
}

// RouteTaskMethods обрабатывает HTTP-запросы для добавления, получения, обновления и удаления задач.
// В зависимости от HTTP-метода маршрутизирует на соответствующую функцию.
func (h *TaskHandler) RouteTaskMethods(w http.ResponseWriter, r *http.Request) {
	// Карта методов HTTP и соответствующих обработчиков
	methods := map[string]func(http.ResponseWriter, *http.Request){
		http.MethodPost:   func(w http.ResponseWriter, r *http.Request) { tasks.AddTaskHandler(w, r, h.db) },
		http.MethodGet:    func(w http.ResponseWriter, r *http.Request) { tasks.GetTaskHandler(w, r, h.db) },
		http.MethodPut:    func(w http.ResponseWriter, r *http.Request) { tasks.UpdateTaskHandler(w, r, h.db) },
		http.MethodDelete: func(w http.ResponseWriter, r *http.Request) { tasks.DeleteTaskHandler(w, r, h.db) },
	}

	// Выбор и вызов обработчика на основе HTTP-метода
	if handlerFunc, exists := methods[r.Method]; exists {
		handlerFunc(w, r)
	} else {
		// Ответ с ошибкой 405, если метод не поддерживается
		tasks.RespondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}
