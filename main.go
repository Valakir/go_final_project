package main

import (
	"go_final_project/dates"
	"log"
	"net/http"
	"os"

	"go_final_project/database"
)

func initDB() {
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
	http.HandleFunc("/api/nextdate", dates.ApiNextDateHandler)

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
