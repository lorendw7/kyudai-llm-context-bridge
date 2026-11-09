package main

import (
	"errors"
	"strings"
	"time"
)

type Session struct {
	ID          string
	LabID       string
	OwnerUserID string
	PolicyID    string
	Revision    uint64
	CreatedAt   time.Time
}

func NewSession(
	id string,
	labID string,
	ownerUserID string,
	policyID string,
	createdAt time.Time,
) (Session, error) {
	if strings.TrimSpace(id) == "" {
		return Session{}, errors.New("session ID is required")
	}

	if strings.TrimSpace(labID) == "" {
		return Session{}, errors.New("lab ID is required")
	}

	if strings.TrimSpace(ownerUserID) == "" {
		return Session{}, errors.New("owner user ID is required")
	}

	if strings.TrimSpace(policyID) == "" {
		return Session{}, errors.New("policy ID is required")
	}

	if createdAt.IsZero() {
		return Session{}, errors.New("creation time is required")
	}

	return Session{
		ID:          id,
		LabID:       labID,
		OwnerUserID: ownerUserID,
		PolicyID:    policyID,
		Revision:    1,
		CreatedAt:   createdAt,
	}, nil
}

func (s Session) OwnedBy(labID string, userID string) bool {
	return s.LabID == labID && s.OwnerUserID == userID
}

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type Message struct {
	ID        string
	SessionID string
	UserID    string
	Role      MessageRole
	Content   string
	CreatedAt time.Time
}

func NewMessage(
	id string,
	sessionID string,
	userID string,
	role MessageRole,
	content string,
	createdAt time.Time,
) (Message, error) {
	if strings.TrimSpace(id) == "" {
		return Message{}, errors.New("message ID is required")
	}

	if strings.TrimSpace(sessionID) == "" {
		return Message{}, errors.New("session ID is required")
	}

	if strings.TrimSpace(content) == "" {
		return Message{}, errors.New("message content is required")
	}

	if createdAt.IsZero() {
		return Message{}, errors.New("creation time is required")
	}

	switch role {
	case MessageRoleUser:
		if strings.TrimSpace(userID) == "" {
			return Message{}, errors.New("user ID is required for user messages")
		}

	case MessageRoleAssistant:
		if strings.TrimSpace(userID) != "" {
			return Message{}, errors.New("assistant messages must not have a user ID")
		}

	default:
		return Message{}, errors.New("unsupported message role")
	}

	return Message{
		ID:        id,
		SessionID: sessionID,
		UserID:    userID,
		Role:      role,
		Content:   content,
		CreatedAt: createdAt,
	}, nil
}
