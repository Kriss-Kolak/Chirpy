package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSimpleHashing(t *testing.T) {
	testPassword := "mySecret"
	hashed, err := HashPassword(testPassword)
	if err != nil || hashed == "" {
		t.Errorf("Encountered %s", err)
	}
}

func TestHashPasswordComparison(t *testing.T) {
	testPassword := "mySecret"
	hashed, err := HashPassword(testPassword)
	if err != nil || hashed == "" {
		t.Errorf("Encountered %v", err)
	}
	valid, err := CheckPasswordHash(testPassword, hashed)
	if err != nil || valid == false {
		t.Errorf("Encountered %v", err)
	}
}

func TestTokenCreation(t *testing.T) {
	var timeDuration time.Duration = 1 * time.Second // 1s
	var secret string = "aj43nfsdf9sdf0sdf9809sd0fsdf0sdf"

	testUserID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Errorf("Encountered %v", err)
	}

	userUUID, err := MakeJWT(testUserID, secret, timeDuration)
	if err != nil || userUUID == "" {
		t.Errorf("Encountered %v", err)
	}

}

func TestTokenInvalidDuration(t *testing.T) {
	var timeDuration time.Duration = 1 * time.Second // 1s
	var secret string = "aj43nfsdf9sdf0sdf9809sd0fsdf0sdf"

	testUserID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Errorf("Encountered %v", err)
	}

	token, err := MakeJWT(testUserID, secret, timeDuration)
	if err != nil || token == "" {
		t.Errorf("Encountered %v", err)
	}

	time.Sleep(2 * time.Second)

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("Encountered %v", err)
	}

}

func TestTokenValidDuration(t *testing.T) {
	var timeDuration time.Duration = 1 * time.Second // 1s
	var secret string = "aj43nfsdf9sdf0sdf9809sd0fsdf0sdf"

	testUserID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Errorf("Encountered %v", err)
	}

	token, err := MakeJWT(testUserID, secret, timeDuration)
	if err != nil || token == "" {
		t.Errorf("Encountered %v", err)
	}

	user_id, err := ValidateJWT(token, secret)
	if err != nil || user_id != testUserID {
		t.Errorf("Encountered %v", err)
	}

}
