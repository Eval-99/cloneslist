package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	params := &argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  16,
		KeyLength:   32,
	}

	hashedPass, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}

	return hashedPass, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	isValid, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return isValid, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSig, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return jwtSig, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	user_id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	parsedUser, err := uuid.Parse(user_id)
	if err != nil {
		return uuid.Nil, err
	}

	return parsedUser, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header missing")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Authorization header must use Bearer scheme")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return "", errors.New("Bearer token is empty")
	}
	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header missing")
	}
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", errors.New("Authorization header must use ApiKey scheme")
	}
	key := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey "))
	if key == "" {
		return "", errors.New("ApiKey key is empty")
	}
	return key, nil
}

func MakeRefreshToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
