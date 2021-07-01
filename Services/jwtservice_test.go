package services

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerateToken(t *testing.T) {
	//Generate token based on sample claim
	sample_token := NewJWTService().GenerateToken("sample@email.com")

	//Try to validate the token
	token_to_check, err := NewJWTService().ValidateToken(sample_token)

	if err != nil {
		t.Fatal(err)
	}

	//Token should be valid
	if !token_to_check.Valid {
		t.Fatalf(`Invalid Token`)
	}
}

func TestValidateToken(t *testing.T) {
	t.Run("InvalidKey", func(t *testing.T) {
		//Try to guess the secretkey
		mock_jwtService := &jwtService{
			secretKey: "invalidKey",
			issuer:    "thinkbridgeIdProvider",
		}

		//Generate token based on sample claim
		sample_token := mock_jwtService.GenerateToken("sample@email.com")

		//Try to validate the token
		_, err := NewJWTService().ValidateToken(sample_token)

		//Validation should err
		if err == nil {
			t.Fatal(`should've been error for invalid key`)
		}
	})

	t.Run("EmptyKey", func(t *testing.T) {
		//Try to get token without secretkey
		mock_jwtService := &jwtService{
			secretKey: "",
			issuer:    "thinkbridgeIdProvider",
		}

		//Generate token based on sample claim
		sample_token := mock_jwtService.GenerateToken("sample@email.com")

		//Try to validate the token
		_, err := NewJWTService().ValidateToken(sample_token)

		//Validation should err
		if err == nil {
			t.Fatal(`should've been error for invalid key`)
		}
	})

	t.Run("InvalidSigningMethod", func(t *testing.T) {
		mock_jwtService := &jwtService{
			secretKey: "UjgFm344XW",
			issuer:    "thinkbridgeIdProvider",
		}

		//Set sample claims
		claims := &jwtCustomClaims{
			"sample@email.com",
			jwt.StandardClaims{
				Issuer:    mock_jwtService.issuer,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}

		//Create token with invalid signing method
		sample_token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

		//Generate encoded token using the secret signing key
		output, _ := sample_token.SignedString([]byte(mock_jwtService.secretKey))

		//Try to validate the token
		_, err := NewJWTService().ValidateToken(output)

		//Validation should err
		if err == nil {
			t.Fatal(`should've been error for unexpected signing method key`)
		}
	})
}
