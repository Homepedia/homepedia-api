package service

import "testing"

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error("HashPassword() should not return an error")
	}
	if hash == "" {
		t.Error("HashPassword() should return a hash")
	}

	// Test if the hash is different from the password
	verifyPassword := VerifyPassword(password, hash)

	if !verifyPassword {
		t.Error("VerifyPassword() should return true")
	}

}
