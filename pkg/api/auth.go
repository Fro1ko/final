package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const tokenSecret = "final-scheduler-token-secret"

type signinRequest struct {
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type tokenClaims struct {
	PasswordHash string `json:"password_hash"`
	ExpiresAt    int64  `json:"exp"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, signinResponse{Error: "method not allowed"})
		return
	}

	var request signinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, signinResponse{Error: err.Error()})
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if subtle.ConstantTimeCompare([]byte(request.Password), []byte(password)) != 1 {
		writeJSON(w, http.StatusUnauthorized, signinResponse{Error: "Неверный пароль"})
		return
	}

	token, err := makeToken(password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, signinResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, signinResponse{Token: token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password != "" {
			cookie, err := r.Cookie("token")
			if err != nil || !validToken(cookie.Value, password) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func makeToken(password string) (string, error) {
	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	claims, err := json.Marshal(tokenClaims{
		PasswordHash: passwordHash(password),
		ExpiresAt:    time.Now().Add(8 * time.Hour).Unix(),
	})
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(header)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claims)
	signingInput := encodedHeader + "." + encodedClaims
	return signingInput + "." + sign(signingInput), nil
}

func validToken(token string, password string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	signingInput := parts[0] + "." + parts[1]
	expectedSignature := sign(signingInput)
	if subtle.ConstantTimeCompare([]byte(parts[2]), []byte(expectedSignature)) != 1 {
		return false
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var claims tokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return false
	}
	return claims.PasswordHash == passwordHash(password) && claims.ExpiresAt > time.Now().Unix()
}

func sign(value string) string {
	mac := hmac.New(sha256.New, []byte(tokenSecret))
	_, _ = mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func passwordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash[:])
}
