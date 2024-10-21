package dates

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Разбираем исходную дату
	originalDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неправильный формат даты: %w", err)
	}

	// Если правило пустое — возвращаем ошибку
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	// Логика обработки правила repeat
	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("некорректный формат правила d")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неправильное количество дней")
		}

		// Использование AddDate с добавлением days
		for nextDate := originalDate.AddDate(0, 0, days); ; nextDate = nextDate.AddDate(0, 0, days) {
			if nextDate.After(now) {
				return nextDate.Format(DateFormat), nil
			}
		}

	case "y":
		for nextDate := originalDate.AddDate(1, 0, 0); ; nextDate = nextDate.AddDate(1, 0, 0) {
			if nextDate.After(now) {
				if nextDate.Month() == time.February && nextDate.Day() == 29 && !isLeapYear(nextDate.Year()) {
					nextDate = nextDate.AddDate(0, 0, 1) // Переход на 1 марта, если не високосный год
				}
				return nextDate.Format(DateFormat), nil
			}
		}

	default:
		return "", errors.New("неподдерживаемое правило повторения")
	}
}

// Функция проверки на високосный год
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// Обработчик для API
func ApiNextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		http.Error(w, "неправильный формат now", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, nextDate)
}
