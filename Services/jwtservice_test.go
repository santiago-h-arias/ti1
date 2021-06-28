package services

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	sample_token := NewJWTService().GenerateToken("sample@email.com")
	token_to_check, err := NewJWTService().ValidateToken(sample_token)

	if err != nil {
		t.Fatal(err)
	}

	if !token_to_check.Valid {
		t.Fatalf(`Invalid Token`)
	}
}

func TestValidateToken(t *testing.T) {
	t.Run("InvalidKey", func(t *testing.T) {
		mock_jwtService := &jwtService{
			secretKey: "invalidKey",
			issuer:    "thinkbridgeIdProvider",
		}

		sample_token := mock_jwtService.GenerateToken("sample.email.com")
		_, err := NewJWTService().ValidateToken(sample_token)
		if err == nil {
			t.Fatal(`should've been error for invalid key`)
		}
	})

	t.Run("EmptyKey", func(t *testing.T) {
		mock_jwtService := &jwtService{
			secretKey: ".",
			issuer:    "thinkbridgeIdProvider",
		}

		sample_token := mock_jwtService.GenerateToken("sample.email.com")
		_, err := NewJWTService().ValidateToken(sample_token)
		if err == nil {
			t.Fatal(`should've been error for invalid key`)
		}
	})
}
