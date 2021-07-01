package services

import (
	"testing"
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
		sample_token := mock_jwtService.GenerateToken("sample.email.com")

		//Try to validate the token
		_, err := NewJWTService().ValidateToken(sample_token)

		//Validation should err
		if err == nil {
			t.Fatal(`should've been error for invalid key`)
		}
	})
}
