package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
    userID := uuid.New()
    secret := "mysecret"

    token, err := MakeJWT(userID, secret, time.Hour)
    if err != nil {
        t.Fatalf("MakeJWT failed: %v", err)
    }

	gotID, err := ValidateJWT(token, secret)
    if err != nil {
        t.Fatalf("ValidateJWT failed: %v", err)
    }

    if gotID != userID {
        t.Errorf("got %v, want %v", gotID, userID)
    }
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mysecret"

	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}


	_, err = ValidateJWT(token, secret)
	if err == nil {
	    t.Errorf("expected error for expired token, got nil")
	}
}
