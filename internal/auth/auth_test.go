package internal

import (
	"testing"
)

func TestSimpleHashing(t *testing.T) {
	test_password := "mySecret"
	hashed, err := HashPassword(test_password)
	if err != nil || hashed == "" {
		t.Errorf("Encountered %s", err)
	}
}

func TestHashPasswordComparison(t *testing.T) {
	test_password := "mySecret"
	hashed, err := HashPassword(test_password)
	if err != nil || hashed == "" {
		t.Errorf("Encountered %v", err)
	}
	valid, err := CheckPasswordHash(test_password, hashed)
	if err != nil || valid == false {
		t.Errorf("Encountered %v", err)
	}
}
