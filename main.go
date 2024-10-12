package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/database"
)

func initDB() {
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to set up database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	log.Println("Database setup successfully")
}

func main() {

	// Инициализация БД
	initDB()
	// Установим директорию, откуда будут раздаваться файлы
	webDir := "./web"
	// Создание обработчика файлового сервера
	fileServer := http.FileServer(http.Dir(webDir))
	// Определение пути для файлового сервера
	http.Handle("/", fileServer)

	// Получение порта из переменной окружения TODO_PORT или использование дефолтного
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	// Запуск сервера и проверка на ошибки
	log.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
