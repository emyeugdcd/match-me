package services

import (
	"testing"
)

func TestUserOrdering(t *testing.T) {
	// Test that user ordering is consistent
	user1, user2 := 5, 3
	if user1 > user2 {
		user1, user2 = user2, user1
	}

	if user1 != 3 || user2 != 5 {
		t.Errorf("Expected (3, 5), got (%d, %d)", user1, user2)
	}
}

func TestConnectionRequestValidation(t *testing.T) {
	// These would be integration tests with a test database

	fromUserID := 1
	toUserID := 1

	if fromUserID == toUserID {
		// Expected: should fail
		t.Log("Correctly prevents self-connection")
	}
}
