package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, *argon2id.Params, error) {
	result, params, err := argon2id.CheckHash(password, hash)
	if err != nil {
		return result, nil, err
	}

	return result, params, nil
}

func OpportunisticRehashing(params *argon2id.Params, password string) (string, error) {
	if params.Memory < argon2id.DefaultParams.Memory || params.Iterations < argon2id.DefaultParams.Iterations {

		return HashPassword(password)
	}
	return "", nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	headerValue := headers.Get("Authorization")
	fmt.Printf("headerValue:%v\n\n\n", headerValue)
	if headerValue == "" {
		return "", errors.New("Couldn't get Authorization header!")
	}

	tokenString := strings.Replace(headerValue, "Bearer", "", 1)
	tokenString = strings.TrimSpace(tokenString)
	fmt.Printf("tokenString:%v\n\n\n", tokenString)
	return tokenString, nil
}

func GetApiKey(headers http.Header) (string, error) {
	headerValue := headers.Get("Authorization")
	fmt.Printf("headerValue:%v\n\n\n", headerValue)
	if headerValue == "" {
		return "", errors.New("Couldn't get Authorization header!")
	}

	apiKeyString := strings.Replace(headerValue, "ApiKey", "", 1)
	apiKeyString = strings.TrimSpace(apiKeyString)
	fmt.Printf("apiKeyString:%v\n\n\n", apiKeyString)
	return apiKeyString, nil
}

func MakeRefreshToken() (string, error) {
	byteToken := make([]byte, 32)
	_, err := rand.Read(byteToken)
	if err != nil {
		return "", err
	}
	fmt.Print(byteToken)

	return hex.EncodeToString(byteToken), nil
}
