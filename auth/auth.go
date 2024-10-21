package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

var password = os.Getenv("TODO_PASSWORD")

// SignInHandler Обработчик для входа в систему возвращающий JWT токен
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Error: "Некорректный запрос"})
		return
	}

	// Хэш представлен в удобочитаемой строке для сравнения
	hashedPassword := ComputeHash(password)
	log.Printf("Хэш пароля из переменной окружения: %s", hashedPassword)

	// Проверка введенного пароля
	if ComputeHash(loginReq.Password) != hashedPassword {
		json.NewEncoder(w).Encode(AuthResponse{Error: "Неверный пароль"})
		return
	}

	// Генерация JWT токена с хэшем пароля
	token, err := createJWT(hashedPassword)
	if err != nil {
		log.Printf("Ошибка создания токена: %v", err)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Ошибка создания токена"})
		return
	}

	// Установка куки
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(8 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	log.Printf("Установлена кука %v", token)
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

// Создание SHA-256 хэша от строки
func ComputeHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// Создание JWT токена
func createJWT(hash string) (string, error) {
	secretKey := []byte(os.Getenv("TODO_PASSWORD"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	return token.SignedString(secretKey)
}

// Аутентификация пользователя с помощью JWT
func AuthUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if password == "" {

			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}
		// Получение куки
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Printf("Ошибка получения куки: %v", err)
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		secretKey := []byte(password)

		// Парсинг токена
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный метод подписи токена")
			}
			return secretKey, nil
		})
		// Проверка токена
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			currentHash := ComputeHash(password)
			tokenHash, ok := claims["hash"].(string)
			log.Printf("Текущий хэш: %s, Хэш в токене: %s", currentHash, tokenHash)
			log.Printf("Пароль: %s", password)        // Log the password
			log.Printf("Хэш пароля: %s", currentHash) // Log the hash of the password

			if !ok || tokenHash != currentHash {
				log.Print("Хэш пароля не совпадает")
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}
		} else {
			log.Printf("Ошибка проверки токена: %v", err)
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
