package main

import (
	"testing"
	"time"
)

func TestNewSessionCreatesOwnedSession(t *testing.T) {
	createdAt := time.Date(2026, time.November, 8, 9, 30, 0, 0, time.UTC)

	session, err := NewSession(
		"session-1",
		"lab-a",
		"user-1",
		"full-history-v1",
		createdAt,
	)
	if err != nil {
		t.Fatalf("NewSession returned an unexpected error: %v", err)
	}

	if session.ID != "session-1" {
		t.Errorf("unexpected session ID: got %q, want %q", session.ID, "session-1")
	}
	if session.LabID != "lab-a" {
		t.Errorf("unexpected lab ID: got %q, want %q", session.LabID, "lab-a")
	}
	if session.OwnerUserID != "user-1" {
		t.Errorf("unexpected owner user ID: got %q, want %q", session.OwnerUserID, "user-1")
	}
	if session.PolicyID != "full-history-v1" {
		t.Errorf("unexpected policy ID: got %q, want %q", session.PolicyID, "full-history-v1")
	}
	if session.Revision != 1 {
		t.Errorf("unexpected initial revision: got %d, want 1", session.Revision)
	}
	if !session.CreatedAt.Equal(createdAt) {
		t.Errorf("unexpected creation time: got %s, want %s", session.CreatedAt, createdAt)
	}
}

func TestSessionOwnedByRequiresLabAndUser(t *testing.T) {
	session, err := NewSession(
		"session-1",
		"lab-a",
		"user-1",
		"full-history-v1",
		time.Date(2026, time.November, 8, 9, 30, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("NewSession returned an unexpected error: %v", err)
	}

	tests := []struct {
		name   string
		labID  string
		userID string
		want   bool
	}{
		{name: "matching owner", labID: "lab-a", userID: "user-1", want: true},
		{name: "same user in another lab", labID: "lab-b", userID: "user-1", want: false},
		{name: "another user in same lab", labID: "lab-a", userID: "user-2", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := session.OwnedBy(tt.labID, tt.userID); got != tt.want {
				t.Errorf("OwnedBy(%q, %q) = %v, want %v", tt.labID, tt.userID, got, tt.want)
			}
		})
	}
}

func TestNewSessionRejectsMissingRequiredFields(t *testing.T) {
	createdAt := time.Date(2026, time.November, 8, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name        string
		id          string
		labID       string
		ownerUserID string
		policyID    string
		createdAt   time.Time
	}{
		{name: "missing session ID", labID: "lab-a", ownerUserID: "user-1", policyID: "full-history-v1", createdAt: createdAt},
		{name: "blank lab ID", id: "session-1", labID: "   ", ownerUserID: "user-1", policyID: "full-history-v1", createdAt: createdAt},
		{name: "missing owner", id: "session-1", labID: "lab-a", policyID: "full-history-v1", createdAt: createdAt},
		{name: "missing policy", id: "session-1", labID: "lab-a", ownerUserID: "user-1", createdAt: createdAt},
		{name: "zero creation time", id: "session-1", labID: "lab-a", ownerUserID: "user-1", policyID: "full-history-v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewSession(tt.id, tt.labID, tt.ownerUserID, tt.policyID, tt.createdAt); err == nil {
				t.Fatal("NewSession returned nil error for invalid input")
			}
		})
	}
}

func TestNewMessageCreatesUserAndAssistantMessages(t *testing.T) {
	createdAt := time.Date(2026, time.November, 8, 9, 35, 0, 0, time.UTC)

	tests := []struct {
		name       string
		userID     string
		role       MessageRole
		content    string
		wantUserID string
		wantRole   MessageRole
	}{
		{
			name:       "user source message",
			userID:     "user-1",
			role:       MessageRoleUser,
			content:    "I prefer tea.",
			wantUserID: "user-1",
			wantRole:   MessageRoleUser,
		},
		{
			name:     "assistant output message",
			role:     MessageRoleAssistant,
			content:  "Understood.",
			wantRole: MessageRoleAssistant,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := NewMessage(
				"message-1",
				"session-1",
				tt.userID,
				tt.role,
				tt.content,
				createdAt,
			)
			if err != nil {
				t.Fatalf("NewMessage returned an unexpected error: %v", err)
			}

			if message.ID != "message-1" {
				t.Errorf("unexpected message ID: got %q, want %q", message.ID, "message-1")
			}
			if message.SessionID != "session-1" {
				t.Errorf("unexpected session ID: got %q, want %q", message.SessionID, "session-1")
			}
			if message.UserID != tt.wantUserID {
				t.Errorf("unexpected user ID: got %q, want %q", message.UserID, tt.wantUserID)
			}
			if message.Role != tt.wantRole {
				t.Errorf("unexpected role: got %q, want %q", message.Role, tt.wantRole)
			}
			if message.Content != tt.content {
				t.Errorf("unexpected content: got %q, want %q", message.Content, tt.content)
			}
			if !message.CreatedAt.Equal(createdAt) {
				t.Errorf("unexpected creation time: got %s, want %s", message.CreatedAt, createdAt)
			}
		})
	}
}

func TestNewMessageRejectsInvalidInput(t *testing.T) {
	createdAt := time.Date(2026, time.November, 8, 9, 35, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        string
		sessionID string
		userID    string
		role      MessageRole
		content   string
		createdAt time.Time
	}{
		{name: "missing message ID", sessionID: "session-1", userID: "user-1", role: MessageRoleUser, content: "hello", createdAt: createdAt},
		{name: "missing session ID", id: "message-1", userID: "user-1", role: MessageRoleUser, content: "hello", createdAt: createdAt},
		{name: "missing user ID for user role", id: "message-1", sessionID: "session-1", role: MessageRoleUser, content: "hello", createdAt: createdAt},
		{name: "user ID supplied for assistant role", id: "message-1", sessionID: "session-1", userID: "user-1", role: MessageRoleAssistant, content: "hello", createdAt: createdAt},
		{name: "unknown role", id: "message-1", sessionID: "session-1", userID: "user-1", role: MessageRole("system"), content: "hello", createdAt: createdAt},
		{name: "blank content", id: "message-1", sessionID: "session-1", userID: "user-1", role: MessageRoleUser, content: "   ", createdAt: createdAt},
		{name: "zero creation time", id: "message-1", sessionID: "session-1", userID: "user-1", role: MessageRoleUser, content: "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewMessage(tt.id, tt.sessionID, tt.userID, tt.role, tt.content, tt.createdAt); err == nil {
				t.Fatal("NewMessage returned nil error for invalid input")
			}
		})
	}
}

func TestMessagesKeepTheirSessionScope(t *testing.T) {
	createdAt := time.Date(2026, time.November, 8, 9, 35, 0, 0, time.UTC)

	first, err := NewMessage("message-1", "session-a", "user-1", MessageRoleUser, "first", createdAt)
	if err != nil {
		t.Fatalf("first NewMessage returned an unexpected error: %v", err)
	}
	second, err := NewMessage("message-2", "session-b", "user-1", MessageRoleUser, "second", createdAt)
	if err != nil {
		t.Fatalf("second NewMessage returned an unexpected error: %v", err)
	}

	if first.SessionID == second.SessionID {
		t.Fatalf("messages from separate sessions share scope %q", first.SessionID)
	}
}
