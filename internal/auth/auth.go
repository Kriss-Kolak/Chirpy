package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	valid, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return valid, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	NowDate := time.Time.UTC(time.Now())
	ExpiredDate := NowDate.Add(expiresIn)

	NewClaim := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(NowDate),
		ExpiresAt: jwt.NewNumericDate(ExpiredDate),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaim)
	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	NewClaim := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &NewClaim, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	if time.Since(NewClaim.ExpiresAt.Time) > 0 {
		return uuid.UUID{}, errors.New("token has expired")
	}

	hash, err := uuid.Parse(NewClaim.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return hash, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("something went wrong")
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("something went wrong")
	}

	userToken := strings.TrimPrefix(token, "Bearer ")
	userToken = strings.TrimSpace(userToken)
	return userToken, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	APIkey := headers.Get("Authorization")
	if APIkey == "" {
		return "", errors.New("something went wrong")
	}
	if !strings.HasPrefix(APIkey, "ApiKey ") {
		return "", errors.New("something went wrong")
	}

	APIkey = strings.TrimPrefix(APIkey, "ApiKey ")
	APIkey = strings.TrimSpace(APIkey)
	return APIkey, nil
}
