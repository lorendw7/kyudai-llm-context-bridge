package main

import (
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	beforeCreation := time.Now()

	message := NewMessage(
		"message-1",
		"session-1",
		"user-1",
		"user",
		"I prefer tea.",
	)

	afterCreation := time.Now()

	if message.ID != "message-1" {
		t.Errorf("unexpected message ID: got %q, want %q", message.ID, "message-1")
	}

	if message.SessionID != "session-1" {
		t.Errorf("unexpected session ID: got %q, want %q", message.SessionID, "session-1")
	}

	if message.UserID != "user-1" {
		t.Errorf("unexpected user ID: got %q, want %q", message.UserID, "user-1")
	}

	if message.Role != "user" {
		t.Errorf("unexpected role: got %q, want %q", message.Role, "user")
	}

	if message.Content != "I prefer tea." {
		t.Errorf("unexpected content: got %q, want %q", message.Content, "I prefer tea.")
	}

	if message.CreatedAt.Before(beforeCreation) || message.CreatedAt.After(afterCreation) {
		t.Errorf(
			"creation time %s is outside the call interval [%s, %s]",
			message.CreatedAt,
			beforeCreation,
			afterCreation,
		)
	}
}
