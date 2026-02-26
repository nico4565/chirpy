package auth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, _, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	cases := []struct {
		testCase    string
		tokenSecret string
		expiration  time.Duration
	}{
		{
			testCase:    "happy path",
			tokenSecret: "secret-a",
			expiration:  time.Hour,
		},
		{
			testCase:    "wrong secret",
			tokenSecret: "secret-b",
			expiration:  time.Hour,
		},
		{
			testCase:    "token expired",
			tokenSecret: "secret-a",
			expiration:  -time.Hour,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", c.testCase), func(t *testing.T) {
			id := uuid.New()
			token, err := MakeJWT(id, "secret-a", c.expiration)
			if err != nil {
				t.Fatalf("failed to make jwt: %v", err)
			}

			switch testCase := c.testCase; testCase {

			case "happy path":
				validatedID, err := ValidateJWT(token, c.tokenSecret)
				if err != nil {
					t.Fatalf("failed to validate jwt: %v", err)
				}

				if validatedID != id {
					t.Errorf("expected %v, got %v", id, validatedID)
				}
			case "wrong secret":
				_, err := ValidateJWT(token, c.tokenSecret)
				if !strings.Contains(err.Error(), "token signature is invalid") {
					t.Fatalf("Validation should have failed, jwt: %v", err)
				}

			case "token expired":
				_, err := ValidateJWT(token, c.tokenSecret)
				if !strings.Contains(err.Error(), "expired") {
					t.Fatalf("Validation should have failed, jwt: %v", err)
				}
				t.Log(err.Error())
			default:
			}

		})
	}
}

func TestGetBearerToken(t *testing.T) {
	headerWithAuth := http.Header{}
	headerWithAuth.Set("Authorization", "Bearer token")

	headerWithoutAuth := http.Header{}

	cases := []struct {
		testCase string
		header   http.Header
	}{
		{
			testCase: "happy path",
			header:   headerWithAuth,
		},
		{
			testCase: "no auth header",
			header:   headerWithoutAuth,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", c.testCase), func(t *testing.T) {

			switch testCase := c.testCase; testCase {

			case "happy path":
				tokenString, _ := GetBearerToken(c.header)
				if tokenString == "" {
					t.Fatal("tokenString empty")
				}
			case "no auth header":
				tokenString, _ := GetBearerToken(c.header)
				if tokenString != "" {
					t.Fatalf("tokenString is not empty as it should, [tokenSting: %s]", tokenString)
				}
			default:
			}

		})
	}
}
