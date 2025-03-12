package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type requestAuth struct {
	Password string `json:"password"`
}

type responseAuth struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func generateToken() string {
	password := os.Getenv("TODO_PASSWORD")
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func validateToken(token string) bool {
	if token == generateToken() {
		return true
	}
	return false
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			valid := validateToken(jwt)

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func handlerAuth(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		errStr, err := json.Marshal(responseAuth{Error: err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(errStr)
		return
	}
	reqA := requestAuth{}
	err = json.Unmarshal(body, &reqA)
	if err != nil {
		errStr, err := json.Marshal(responseAuth{Error: err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(errStr)
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if password != reqA.Password {
		w.WriteHeader(http.StatusUnauthorized)
		errStr, err := json.Marshal(responseAuth{Error: "неверный пароль"})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(errStr)
		return
	} else {
		token := generateToken()
		respStr, err := json.Marshal(responseAuth{Token: token})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respStr)
		return
	}
}
