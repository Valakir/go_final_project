package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Valakir/go_final_project/database"
)

func main() {
	// Инициализация БД
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("Ошибка настройки БД: %v", err)
	}
	defer db.Close()

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
