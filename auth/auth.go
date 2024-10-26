package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

var passwordENV = os.Getenv("TODO_PASSWORD")

// Создание SHA-256 хэша от строки
func ComputeHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Создание JWT токена
func createJWT(hash string) (string, error) {
	secretKey := []byte(passwordENV)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	return token.SignedString(secretKey)
}

// SignInHandler обработчик для входа в систему, возвращающий JWT токен
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Проверяем, что запрос является POST-запросом
	if r.Method != "POST" {
		http.Error(w, "Метод запроса не POST", http.StatusMethodNotAllowed)
		return
	}
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	// Проверяем, что пароль не пустой
	if loginReq.Password == "" {
		http.Error(w, "Пароль не может быть пустым", http.StatusBadRequest)
		return
	}
	var passwordENV = os.Getenv("TODO_PASSWORD")
	// Проверяем, что переменная окружения для пароля не пустая
	if passwordENV == "" {
		http.Error(w, "Ошибка сервера: недоступен пароль для проверки", http.StatusInternalServerError)
		return
	}

	// Сравниваем хэши
	if ComputeHash(loginReq.Password) != ComputeHash(passwordENV) {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	token, err := createJWT(ComputeHash(loginReq.Password))
	if err != nil {
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(8 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

// AuthUser аутентификация пользователя с помощью JWT
func AuthUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if passwordENV == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		secretKey := []byte(passwordENV)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			currentHash := ComputeHash(passwordENV)
			if tokenHash, ok := claims["hash"].(string); !ok || tokenHash != currentHash {
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
