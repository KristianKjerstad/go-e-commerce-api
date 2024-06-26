package auth

import "testing"

func TestJWT(t *testing.T) {

	t.Run("Hash passwords", createJWTTest)

}

func createJWTTest(t *testing.T) {
	secret := "secret-string"
	userID := 52

	token, err := CreateJWT([]byte(secret), userID)
	if err != nil {
		t.Errorf("error creating jwt: %v", err)
	}

	if token == "" {
		t.Errorf("error creating jwt, value was empty")
	}

}
