package auth

import "testing"

func TestHashPassword(t *testing.T) {

	t.Run("Hash passwords", hashPasswordTest)

}

func hashPasswordTest(t *testing.T) {
	password := "secret-password"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hashed == "" {
		t.Errorf("error hashing password, hashed value was empty")
	}

	arePasswordsEqual := ComparePasswords(hashed, []byte(password))
	if !arePasswordsEqual {
		t.Errorf("error comparing passwords.")
	}
}
