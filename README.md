# Приложение TODO

## Описание

Приложение TODO — это веб-приложение, которое позволяет создавать и удалять задачи.
В директории
- `main.go` содержит точку входа в приложение
- `database` содержит БД и функции для работы с ней
- `tasks` содержит классы для работы с задачами
- `dates` содержит функции для работы с датами
- `auth` содержит функции для работы с авторизацией
- `tests` находятся тесты для проверки API
- `web` содержит файлы фронтенда.

## Инструкции
**Запуск в докере:**
- docker build --tag my_app:v1 . 
- docker run -d -p 7540:7540 my_app:v1 
- страница вызова приложения: http://localhost:7540

## Что выполнено со звёздочкой
- TODO_PORT — порт, на котором будет работать приложение
- TODO_DBFILE - путь к БД
- Dockerfile

### Тестирование и отладка
В папке `tests` находятся файлы тестов. 

Запуск одного теста:
go test -run ^TestDB$ ./tests -count=1
go test -run ^TestNextDate$ ./tests -count=1
go test -run ^TestAddTask$ ./tests -count=1
go test -run ^TestTasks$ ./tests -count=1
go test -run ^TestEditTask$ ./tests -count=1
go test -run ^TestDone$ ./tests -count=1
go test -run ^TestDelTask$ ./tests -count=1

Запуск всех тестов последовательно:
go test ./tests -count=1