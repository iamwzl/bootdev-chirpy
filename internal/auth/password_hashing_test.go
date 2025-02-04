package auth

import (
    "testing"
)

func TestHashPassword(t *testing.T) {
    password := "TestyTestTester123"
    hash, err := HashPassword(password)
    if err != nil {
        t.Errorf("HashPassword failed: %v", err)
    }
    if hash == password {
        t.Error("Hash should not be equal to password")
    }
}

func TestCheckPasswordHash(t *testing.T) {
    password := "AppleOrange@BananaPear123"
    hash, _ := HashPassword(password)

    err := CheckPasswordHash(hash, password)
    if err != nil {
        t.Error("CheckPasswordHash failed with matching password")
    }

    wrongPassword := "smashthestate"
    err = CheckPasswordHash(hash, wrongPassword)
    if err == nil {
        t.Error("CheckPasswordHash should fail with wrong password")
    }
}
